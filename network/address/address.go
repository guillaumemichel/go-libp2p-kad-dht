package address

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type NetworkAddress interface {
	NodeID() NodeID
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
