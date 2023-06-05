package libp2pendpoint

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO: Use sync.Pool to reuse buffers https://pkg.go.dev/sync#Pool

type Libp2pEndpoint struct {
	host host.Host

	// peer filters to be applied before adding peer to peerstore

	writers sync.Pool
	readers sync.Pool
}

func NewMessageEndpoint(host host.Host) *Libp2pEndpoint {
	return &Libp2pEndpoint{
		host:    host,
		writers: sync.Pool{},
		readers: sync.Pool{},
	}
}

func (msgEndpoint *Libp2pEndpoint) AsyncDialAndReport(ctx context.Context, p peer.ID, reportFn endpoint.DialReportFn) {
	go func() {
		ctx, span := internal.StartSpan(ctx, "Libp2pEndpoint.AsyncDialAndReport", trace.WithAttributes(
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

func (msgEndpoint *Libp2pEndpoint) DialPeer(ctx context.Context, p peer.ID) error {
	_, span := internal.StartSpan(ctx, "Libp2pEndpoint.DialPeer", trace.WithAttributes(
		attribute.String("PeerID", p.String()),
	))
	defer span.End()

	if msgEndpoint.host.Network().Connectedness(p) == network.Connected {
		span.AddEvent("Already connected")
		return nil
	}

	pi := peer.AddrInfo{ID: p}
	if err := msgEndpoint.host.Connect(ctx, pi); err != nil {
		span.AddEvent("Connection failed", trace.WithAttributes(
			attribute.String("Error", err.Error()),
		))
		return err
	}
	span.AddEvent("Connection successful")
	return nil
}

func (msgEndpoint *Libp2pEndpoint) MaybeAddToPeerstore(ai peer.AddrInfo, ttl time.Duration) {
	// Don't add addresses for self or our connected peers. We have better ones.
	if ai.ID == msgEndpoint.host.ID() ||
		msgEndpoint.host.Network().Connectedness(ai.ID) == network.Connected {
		return
	}
	msgEndpoint.host.Peerstore().AddAddrs(ai.ID, ai.Addrs, ttl)
}

func (msgEndpoint *Libp2pEndpoint) SendRequest(ctx context.Context, p peer.ID, req *pb.Message, proto protocol.ID) (*pb.Message, error) {
	ctx, span := internal.StartSpan(ctx, "Libp2pEndpoint.SendRequest", trace.WithAttributes(
		attribute.String("PeerID", p.String()),
	))
	defer span.End()

	s, err := msgEndpoint.host.NewStream(ctx, p, proto)
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

func (msgEndpoint *Libp2pEndpoint) Connectedness(p peer.ID) network.Connectedness {
	return msgEndpoint.host.Network().Connectedness(p)
}

func (msgEndpoint *Libp2pEndpoint) PeerInfo(p peer.ID) peer.AddrInfo {
	return msgEndpoint.host.Peerstore().PeerInfo(p)
}

func (e *Libp2pEndpoint) KadID() key.KadKey {
	return key.PeerKadID(e.host.ID())
}
