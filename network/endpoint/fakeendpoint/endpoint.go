package fakeendpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/events/dispatch"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type FakeEndpoint struct {
	self      address.NodeID
	peerstore map[address.NodeID]peer.AddrInfo

	dispatcher dispatch.Dispatcher

	connStatus map[address.NodeID]network.Connectedness
}

func NewFakeEndpoint() *FakeEndpoint {
	return &FakeEndpoint{
		peerstore: make(map[address.NodeID]peer.AddrInfo),
	}
}

func (e *FakeEndpoint) AsyncDialAndReport(ctx context.Context, id address.NodeID,
	reportFn endpoint.DialReportFn) {

}

func (e *FakeEndpoint) DialPeer(ctx context.Context, id address.NodeID) error {
	return nil
}

func (e *FakeEndpoint) MaybeAddToPeerstore(na address.NetworkAddress, ttl time.Duration) {

}

func (e *FakeEndpoint) SendRequest(ctx context.Context, id address.NodeID,
	req message.MinKadRequestMessage, resp message.MinKadResponseMessage,
	proto protocol.ID) error {

	return nil
}

// Peerstore functions
func (e *FakeEndpoint) Connectedness(id address.NodeID) network.Connectedness {
	if _, ok := e.connStatus[id]; !ok {
		return network.NotConnected
	}
	return e.connStatus[id]
}

func (e *FakeEndpoint) PeerInfo(id address.NodeID) peer.AddrInfo {
	if ai, ok := e.peerstore[id]; ok {
		return ai
	}
	return peer.AddrInfo{}
}

func (e *FakeEndpoint) KadID() key.KadKey {
	return address.KadID(e.self)
}
