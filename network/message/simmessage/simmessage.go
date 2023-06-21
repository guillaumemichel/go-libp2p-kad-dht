package simmessage

import (
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

type SimMessage struct {
	target      *key.KadKey
	closerPeers []address.NodeID
}

func NewSimRequest(target key.KadKey) *SimMessage {
	return &SimMessage{
		target: &target,
	}
}

func NewSimResponse(closerPeers []address.NodeID) *SimMessage {
	return &SimMessage{
		closerPeers: closerPeers,
	}
}

func (m *SimMessage) Target() *key.KadKey {
	return m.target
}

func (m *SimMessage) CloserNodes() []address.NetworkAddress {
	if m.closerPeers == nil {
		return nil
	}
	nas := make([]address.NetworkAddress, len(m.closerPeers))
	for i, peer := range m.closerPeers {
		nas[i] = peer
	}
	return nas
}
