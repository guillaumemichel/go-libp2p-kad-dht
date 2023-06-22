package fakeendpoint

import (
	"context"
	"time"

	ba "github.com/libp2p/go-libp2p-kad-dht/events/action/basicaction"
	sd "github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/libp2p/go-libp2p/core/network"
)

type FakeEndpoint struct {
	self address.NodeID

	peerstore  map[string]address.NetworkAddress
	connStatus map[string]network.Connectedness

	dispatcher *sd.SimpleDispatcher
}

var _ endpoint.Endpoint = (*FakeEndpoint)(nil)

func NewFakeEndpoint(self address.NodeID, dispatcher *sd.SimpleDispatcher) *FakeEndpoint {
	return &FakeEndpoint{
		self:       self,
		peerstore:  make(map[string]address.NetworkAddress),
		connStatus: make(map[string]network.Connectedness),

		dispatcher: dispatcher,
	}
}

func (e *FakeEndpoint) DialPeer(ctx context.Context, id address.NodeID) error {
	_, span := util.StartSpan(ctx, "DialPeer",
		trace.WithAttributes(attribute.String("id", id.String())),
	)
	defer span.End()

	status, ok := e.connStatus[id.String()]

	if ok {
		switch status {
		case network.Connected:
			return nil
		case network.CanConnect:
			e.connStatus[id.String()] = network.Connected
			return nil
		}
	}
	span.RecordError(endpoint.ErrUnknownPeer)
	return endpoint.ErrUnknownPeer
}

// MaybeAddToPeerstore adds the given address to the peerstore. FakeEndpoint
// doesn't take into account the ttl.
func (e *FakeEndpoint) MaybeAddToPeerstore(ctx context.Context, na address.NetworkAddress, ttl time.Duration) error {
	strNodeID := na.NodeID().String()
	_, span := util.StartSpan(ctx, "MaybeAddToPeerstore",
		trace.WithAttributes(attribute.String("self", e.self.String())),
		trace.WithAttributes(attribute.String("id", strNodeID)),
	)
	defer span.End()

	if _, ok := e.peerstore[strNodeID]; !ok {
		e.peerstore[strNodeID] = na
	}
	if _, ok := e.connStatus[strNodeID]; !ok {
		e.connStatus[strNodeID] = network.CanConnect
	}
	return nil
}

func (e *FakeEndpoint) SendRequestHandleResponse(ctx context.Context,
	protoID address.ProtocolID, id address.NodeID, msg message.MinKadMessage,
	resp message.MinKadMessage, handleResp endpoint.ResponseHandlerFn) {

	ctx, span := util.StartSpan(ctx, "SendRequestHandleResponse",
		trace.WithAttributes(attribute.Stringer("id", id)),
	)
	defer span.End()

	if e.DialPeer(ctx, id) != nil {
		return
	}

	replyFn := func(msg message.MinKadResponseMessage) {
		e.dispatcher.DispatchTo(ctx, e.self, ba.BasicAction(func(ctx context.Context) {
			handleResp(ctx, msg, nil)
		}))
	}

	req := msg
	remoteServ := e.dispatcher.Server(id)
	if remoteServ == nil {
		return
	}
	action := ba.BasicAction(func(ctx context.Context) {
		remoteServ.HandleFindNodeRequest(ctx, e.self, req, replyFn)
	})
	e.dispatcher.DispatchTo(ctx, id, action)
}

// Peerstore functions
func (e *FakeEndpoint) Connectedness(id address.NodeID) network.Connectedness {
	if s, ok := e.connStatus[id.String()]; !ok {
		return network.NotConnected
	} else {
		return s
	}
}

func (e *FakeEndpoint) NetworkAddress(id address.NodeID) (address.NetworkAddress, error) {
	if ai, ok := e.peerstore[id.String()]; ok {
		return ai, nil
	}
	return peerid.PeerID{}, endpoint.ErrUnknownPeer
}

func (e *FakeEndpoint) KadKey() key.KadKey {
	return e.self.Key()
}
