package network

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO: Use sync.Pool to reuse buffers https://pkg.go.dev/sync#Pool

type DialReportFn func(context.Context, bool)

type MessageEndpoint struct {
	Host host.Host

	// peer filters to be applied before adding peer to peerstore

	writers sync.Pool
	readers sync.Pool
}

func NewMessageEndpoint(host host.Host) *MessageEndpoint {
	return &MessageEndpoint{
		Host:    host,
		writers: sync.Pool{},
		readers: sync.Pool{},
	}
}

func (msgEndpoint *MessageEndpoint) AsyncDialAndReport(ctx context.Context, p peer.ID, reportFn DialReportFn) {
	go func() {
		ctx, span := internal.StartSpan(ctx, "MessageEndpoint.AsyncDialAndReport", trace.WithAttributes(
			attribute.String("PeerID", p.String()),
		))
		defer span.End()

		success := true
		if err := msgEndpoint.DialPeer(ctx, p); err != nil {
			span.AddEvent("dial failed", trace.WithAttributes(
				attribute.String("Error", err.Error()),
			))
			success = false
		} else {
			span.AddEvent("dial successful")
		}

		// report dial result where it is needed
		reportFn(ctx, success)
	}()
}

func (msgEndpoint *MessageEndpoint) DialPeer(ctx context.Context, p peer.ID) error {
	_, span := internal.StartSpan(ctx, "MessageEndpoint.DialPeer", trace.WithAttributes(
		attribute.String("PeerID", p.String()),
	))
	defer span.End()

	if msgEndpoint.Host.Network().Connectedness(p) == network.Connected {
		span.AddEvent("Already connected")
		return nil
	}

	pi := peer.AddrInfo{ID: p}
	if err := msgEndpoint.Host.Connect(ctx, pi); err != nil {
		span.AddEvent("Connection failed", trace.WithAttributes(
			attribute.String("Error", err.Error()),
		))
		return err
	}
	span.AddEvent("Connection successful")
	return nil
}

func (msgEndpoint *MessageEndpoint) MaybeAddToPeerstore(ai peer.AddrInfo, ttl time.Duration) {
	// Don't add addresses for self or our connected peers. We have better ones.
	if ai.ID == msgEndpoint.Host.ID() ||
		msgEndpoint.Host.Network().Connectedness(ai.ID) == network.Connected {
		return
	}
	msgEndpoint.Host.Peerstore().AddAddrs(ai.ID, ai.Addrs, ttl)
}

func (msgEndpoint *MessageEndpoint) SendRequest(ctx context.Context, p peer.ID, req *pb.Message, proto protocol.ID) (*pb.Message, error) {
	ctx, span := internal.StartSpan(ctx, "MessageEndpoint.SendRequest", trace.WithAttributes(
		attribute.String("PeerID", p.String()),
	))
	defer span.End()

	s, err := msgEndpoint.Host.NewStream(ctx, p, proto)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer s.Close()

	err = WriteMsg(s, req)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := &pb.Message{}
	err = ReadMsg(s, resp)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return resp, nil
}

func WriteMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	w := pbio.NewDelimitedWriter(s)
	return w.WriteMsg(msg)
}

func ReadMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	r := pbio.NewDelimitedReader(s, network.MessageSizeMax)
	return r.ReadMsg(msg)
}
