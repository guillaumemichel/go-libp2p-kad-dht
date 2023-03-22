package simplert

import (
	"sort"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	BUCKET_SIZE = 20
)

type RoutingTable interface {
	AddPeer(peer.AddrInfo) bool
	RemovePeer(hash.KadKey) bool
	NearestPeers(hash.KadKey, int) []peer.AddrInfo
	Find(hash.KadKey) peer.AddrInfo
}

type peerInfo struct {
	id    peer.AddrInfo
	KadId hash.KadKey
}

type DhtRoutingTable struct {
	self    hash.KadKey
	buckets [][]peerInfo
}

func NewDhtRoutingTable(self hash.KadKey) *DhtRoutingTable {
	rt := DhtRoutingTable{
		self:    self,
		buckets: make([][]peerInfo, 0),
	}
	// define bucket 0
	rt.buckets = append(rt.buckets, make([]peerInfo, 0))
	return &rt
}

func (rt *DhtRoutingTable) BucketIdForKey(kadId hash.KadKey) int {
	return hash.CommonPrefixLength(rt.self, kadId)
}

func (rt *DhtRoutingTable) AddPeer(pi peer.AddrInfo) bool {

	kadId := hash.PeerKadID(pi.ID)
	bid := rt.BucketIdForKey(kadId)

	lastBucketId := len(rt.buckets) - 1

	if bid < lastBucketId {
		// new peer doesn't belong in last bucket
		if len(rt.buckets[bid]) >= BUCKET_SIZE {
			// bucket is full, discard new peer
			return false
		}

		for _, p := range rt.buckets[bid] {
			if p.id.ID == pi.ID {
				// peer already in bucket, discard new peer
				return false
			}
		}
		// add new peer to bucket
		rt.buckets[bid] = append(rt.buckets[bid], peerInfo{pi, kadId})
		return true
	}
	if len(rt.buckets[lastBucketId]) < BUCKET_SIZE {
		// last bucket is not full, add new peer
		rt.buckets[lastBucketId] = append(rt.buckets[lastBucketId], peerInfo{pi, kadId})
		return true
	}
	// last bucket is full, try to split it
	for len(rt.buckets[lastBucketId]) == BUCKET_SIZE {
		// farBucket contains peers with a CPL matching lastBucketId
		farBucket := make([]peerInfo, 0)
		// closeBucket contains peers with a CPL higher than lastBucketId
		closeBucket := make([]peerInfo, 0)

		for _, p := range rt.buckets[lastBucketId] {
			if rt.BucketIdForKey(p.KadId) == lastBucketId {
				farBucket = append(farBucket, p)
			} else {
				closeBucket = append(closeBucket, p)
			}
		}
		if len(farBucket) == BUCKET_SIZE {
			// if all peers in the last bucket have the CPL matching this bucket,
			// don't split it and discard the new peer
			return false
		}
		// replace last bucket with farBucket
		rt.buckets[lastBucketId] = farBucket
		// add closeBucket as a new bucket
		rt.buckets = append(rt.buckets, closeBucket)

		lastBucketId++
	}

	return true
}

func (rt *DhtRoutingTable) RemovePeer(kadId hash.KadKey) bool {
	bid := rt.BucketIdForKey(kadId)
	for i, p := range rt.buckets[bid] {
		if p.KadId == kadId {
			// remove peer from bucket
			rt.buckets[bid][i] = rt.buckets[bid][len(rt.buckets[bid])-1]
			rt.buckets[bid] = rt.buckets[bid][:len(rt.buckets[bid])-1]
			return true
		}
	}
	// peer not found in the routing table
	return false
}

func (rt *DhtRoutingTable) Find(kadId hash.KadKey) peer.AddrInfo {
	bid := rt.BucketIdForKey(kadId)
	if bid >= len(rt.buckets) {
		bid = len(rt.buckets) - 1
	}
	for _, p := range rt.buckets[bid] {
		if p.KadId == kadId {
			return p.id
		}
	}
	return peer.AddrInfo{}
}

// TODO: not working as expected
// returns min(n, bucketSize) peers from the bucket matching the given key
func (rt *DhtRoutingTable) NearestPeers(kadId hash.KadKey, n int) []peer.AddrInfo {
	bid := rt.BucketIdForKey(kadId)
	if bid >= len(rt.buckets) {
		bid = len(rt.buckets) - 1
	}
	peers := make([]peerInfo, len(rt.buckets[bid]))
	copy(peers, rt.buckets[bid])
	sort.SliceStable(peers, func(i, j int) bool {
		for k := 0; k < hash.Keysize; k++ {
			distI := peers[i].KadId[k] ^ kadId[k]
			distJ := peers[j].KadId[k] ^ kadId[k]
			if distI != distJ {
				return distI < distJ
			}
		}
		return false
	})
	ais := make([]peer.AddrInfo, min(n, len(peers)))
	for i := 0; i < min(n, len(peers)); i++ {
		ais[i] = peers[i].id
	}

	return ais
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
