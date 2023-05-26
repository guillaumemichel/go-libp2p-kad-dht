package routing

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type Routing interface {
	FindClosestPeers(p peer.ID) ([]peer.ID, error)
	// FindProviders(cid.Cid) ([]peer.ID, error)

	NClosestPeers() int
}
