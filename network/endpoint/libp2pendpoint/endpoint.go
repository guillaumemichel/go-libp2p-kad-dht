package libp2pendpoint

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	laddr "github.com/libp2p/go-libp2p-kad-dht/network/address/libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO: Use sync.Pool to reuse buffers https://pkg.go.dev/sync#Pool

type Libp2pEndpoint struct {
	host    host.Host
	protoID protocol.ID

	// peer filters to be applied before adding peer to peerstore

	writers sync.Pool
	readers sync.Pool
}

func NewMessageEndpoint(host host.Host, proto protocol.ID) *Libp2pEndpoint {
	return &Libp2pEndpoint{
		host:    host,
		protoID: proto,
		writers: sync.Pool{},
		readers: sync.Pool{},
	}
}

func getPeerID(id address.NodeID) peer.ID {
	if p, ok := id.(peer.ID); ok {
		return p
	}
	panic("invalid peer id")
}

func (msgEndpoint *Libp2pEndpoint) AsyncDialAndReport(ctx context.Context, id address.NodeID, reportFn endpoint.DialReportFn) {
	p := getPeerID(id)
	go func() {
		ctx, span := util.StartSpan(ctx, "Libp2pEndpoint.AsyncDialAndReport", trace.WithAttributes(
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

func (msgEndpoint *Libp2pEndpoint) DialPeer(ctx context.Context, id address.NodeID) error {
	p := getPeerID(id)

	_, span := util.StartSpan(ctx, "Libp2pEndpoint.DialPeer", trace.WithAttributes(
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

func (msgEndpoint *Libp2pEndpoint) MaybeAddToPeerstore(na address.NetworkAddress, ttl time.Duration) {
	ai, ok := na.(laddr.Libp2pAddr)
	if !ok {
		panic("invalid peer.AddrInfo")
	}

	// Don't add addresses for self or our connected peers. We have better ones.
	if ai.ID == msgEndpoint.host.ID() ||
		msgEndpoint.host.Network().Connectedness(ai.ID) == network.Connected {
		return
	}
	msgEndpoint.host.Peerstore().AddAddrs(ai.ID, ai.Addrs, ttl)
}

func (msgEndpoint *Libp2pEndpoint) SendRequest(ctx context.Context, id address.NodeID, req message.MinKadRequestMessage,
	resp message.MinKadResponseMessage) error {

	protoReq, ok := req.(message.ProtoKadRequestMessage)
	if !ok {
		panic("Libp2pEndpoint requires ProtoKadRequestMessage")
	}
	protoResp, ok := resp.(message.ProtoKadResponseMessage)
	if !ok {
		panic("Libp2pEndpoint requires ProtoKadResponseMessage")
	}

	p := getPeerID(id)

	ctx, span := util.StartSpan(ctx, "Libp2pEndpoint.SendRequest", trace.WithAttributes(
		attribute.String("PeerID", p.String()),
	))
	defer span.End()

	s, err := msgEndpoint.host.NewStream(ctx, p, msgEndpoint.protoID)
	if err != nil {
		span.RecordError(err)
		return err
	}
	defer s.Close()

	err = WriteMsg(s, protoReq)
	if err != nil {
		span.RecordError(err)
		return err
	}

	err = ReadMsg(s, protoResp)
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (msgEndpoint *Libp2pEndpoint) Connectedness(id address.NodeID) network.Connectedness {
	p := getPeerID(id)
	return msgEndpoint.host.Network().Connectedness(p)
}

func (msgEndpoint *Libp2pEndpoint) PeerInfo(id address.NodeID) peer.AddrInfo {
	p := getPeerID(id)
	return msgEndpoint.host.Peerstore().PeerInfo(p)
}

func (e *Libp2pEndpoint) KadID() key.KadKey {
	return key.PeerKadID(e.host.ID())
}
