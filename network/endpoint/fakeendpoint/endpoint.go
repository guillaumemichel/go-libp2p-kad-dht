package fakeendpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dispatch"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type FakeEndpoint struct {
	self      peer.ID
	peerstore map[peer.ID]peer.AddrInfo

	dispatcher dispatch.Dispatcher

	connStatus map[peer.ID]network.Connectedness
}

func NewFakeEndpoint() *FakeEndpoint {
	return &FakeEndpoint{
		peerstore: make(map[peer.ID]peer.AddrInfo),
	}
}

func (e *FakeEndpoint) AsyncDialAndReport(ctx context.Context, p peer.ID, reportFn endpoint.DialReportFn) {

}

func (e *FakeEndpoint) DialPeer(ctx context.Context, p peer.ID) error {
	return nil
}

func (e *FakeEndpoint) MaybeAddToPeerstore(ai peer.AddrInfo, ttl time.Duration) {

}

func (e *FakeEndpoint) SendRequest(ctx context.Context, p peer.ID, req *pb.Message, proto protocol.ID) (*pb.Message, error) {
	return nil, nil
}

// Peerstore functions
func (e *FakeEndpoint) Connectedness(p peer.ID) network.Connectedness {
	if _, ok := e.connStatus[p]; !ok {
		return network.NotConnected
	}
	return e.connStatus[p]
}

func (e *FakeEndpoint) PeerInfo(p peer.ID) peer.AddrInfo {
	if ai, ok := e.peerstore[p]; ok {
		return ai
	}
	return peer.AddrInfo{}
}

func (e *FakeEndpoint) KadID() key.KadKey {
	return key.PeerKadID(e.self)
}
