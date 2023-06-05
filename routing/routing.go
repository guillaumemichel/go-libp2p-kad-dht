package routing

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

type Routing interface {
	FindClosestPeers(key.KadKey) ([]address.NodeID, error)

	NClosestPeers() int
}
