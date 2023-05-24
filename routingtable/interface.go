package routingtable

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type RoutingTable interface {
	AddPeer(peer.ID) bool
	RemovePeer(key.KadKey) bool
	NearestPeers(key.KadKey, int) []peer.ID
}
