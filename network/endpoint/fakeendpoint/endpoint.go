package fakeendpoint

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	sd "github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type FakeEndpoint struct {
	self address.NodeID

	peerstore  map[address.NodeID]address.NetworkAddress
	connStatus map[address.NodeID]network.Connectedness

	dispatcher *sd.SimpleDispatcher
}

func NewFakeEndpoint(clk *clock.Mock, dispatcher *sd.SimpleDispatcher) *FakeEndpoint {
	return &FakeEndpoint{
		peerstore:  make(map[address.NodeID]address.NetworkAddress),
		connStatus: make(map[address.NodeID]network.Connectedness),

		dispatcher: dispatcher,
	}
}

func (e *FakeEndpoint) AsyncDialAndReport(ctx context.Context, id address.NodeID,
	reportFn endpoint.DialReportFn) {

	reportFn(ctx, e.DialPeer(ctx, id) == nil)
}

func (e *FakeEndpoint) DialPeer(ctx context.Context, id address.NodeID) error {
	status, ok := e.connStatus[id]

	if ok {
		switch status {
		case network.Connected:
			return nil
		case network.CanConnect:
			e.connStatus[id] = network.Connected
			return nil
		case network.CannotConnect:
			return endpoint.ErrCannotConnect
		}
	}
	return endpoint.ErrUnknownPeer
}

// MaybeAddToPeerstore adds the given address to the peerstore. FakeEndpoint
// doesn't take into account the ttl.
func (e *FakeEndpoint) MaybeAddToPeerstore(na address.NetworkAddress, ttl time.Duration) error {
	if _, ok := e.peerstore[address.ID(na)]; !ok {
		e.peerstore[address.ID(na)] = na
	}
	if _, ok := e.connStatus[address.ID(na)]; !ok {
		e.connStatus[address.ID(na)] = network.CanConnect
	}
	return nil
}

func (e *FakeEndpoint) SendRequest(ctx context.Context, id address.NodeID,
	req message.MinKadMessage, resp message.MinKadMessage) error {

	return nil
}

func (e *FakeEndpoint) SendRequestHandleResponse(ctx context.Context, id address.NodeID,
	msg message.MinKadMessage, handleResp endpoint.ResponseHandlerFn) {

	ctx, span := util.StartSpan(ctx, "SendRequestHandleResponse",
		trace.WithAttributes(attribute.Stringer("id", id)),
	)
	defer span.End()

	if e.DialPeer(ctx, id) != nil {
		return
	}

	req := msg.(*ipfskadv1.Message)
	remoteServ := e.dispatcher.Server(id)
	action := func() {
		remoteServ.HandleFindNodeRequest(ctx, e.self, req, handleResp)
	}
	e.dispatcher.DispatchTo(ctx, id, action)
}

// Peerstore functions
func (e *FakeEndpoint) Connectedness(id address.NodeID) network.Connectedness {
	if s, ok := e.connStatus[id]; !ok {
		return network.NotConnected
	} else {
		return s
	}
}

func (e *FakeEndpoint) NetworkAddress(id address.NodeID) (address.NetworkAddress, error) {
	if ai, ok := e.peerstore[id]; ok {
		return ai, nil
	}
	return peer.AddrInfo{}, endpoint.ErrUnknownPeer
}

func (e *FakeEndpoint) KadID() key.KadKey {
	return address.KadID(e.self)
}
