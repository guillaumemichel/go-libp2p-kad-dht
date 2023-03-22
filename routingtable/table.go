package routingtable

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	BUCKET_SIZE = 20
)

type RoutingTable interface {
	AddPeer(peer.AddrInfo) bool
	RemovePeer(hash.KadKey) bool

	
}

type peerInfo struct {
	PeerID peer.ID
	KadId  hash.KadKey
}

type DhtRoutingTable struct {
	buckets [][BUCKET_SIZE]peerInfo
}

func NewDhtRoutingTable() *DhtRoutingTable {
	return &DhtRoutingTable{
		buckets: make([][BUCKET_SIZE]peerInfo, 0),
	}
}
