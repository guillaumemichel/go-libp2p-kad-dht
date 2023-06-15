package address

import (
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type NetworkAddress interface {
}

func ID(na NetworkAddress) NodeID {
	switch na := na.(type) {
	case peer.AddrInfo:
		return na.ID
	case peer.ID:
		return na
	}
	return StringID{}
}

type NodeID interface {
	String() string
}

func KadID(id NodeID) key.KadKey {
	switch id := id.(type) {
	case key.KadKey:
		return id
	case peer.ID:
		return key.PeerKadID(id)
	}

	return key.StringKadID(id.String())
}

type ProtocolID string
