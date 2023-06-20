package simmessage

import (
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/kadid"
)

type SimMessage struct {
	target      kadid.KadID
	closerPeers []kadid.KadID
}

func NewSimRequest(target kadid.KadID) *SimMessage {
	return &SimMessage{
		target: target,
	}
}

func NewSimResponse(closerPeers []kadid.KadID) *SimMessage {
	return &SimMessage{
		closerPeers: closerPeers,
	}
}

func (m *SimMessage) Target() key.KadKey {
	return m.target.KadKey
}

func (m *SimMessage) CloserNodes() []address.NetworkAddress {
	nas := make([]address.NetworkAddress, len(m.closerPeers))
	for i, peer := range m.closerPeers {
		nas[i] = peer
	}
	return nas
}
