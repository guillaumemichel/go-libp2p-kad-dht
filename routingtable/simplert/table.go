package simplert

import (
	"sort"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/peer"
)

type RoutingTable interface {
	AddPeer(peer.AddrInfo) bool
	RemovePeer(hash.KadKey) bool
	NearestPeers(hash.KadKey, int) []peer.AddrInfo
	Find(hash.KadKey) peer.AddrInfo

	BucketSize() int
}

type peerInfo struct {
	id    peer.AddrInfo
	kadId hash.KadKey
}

type DhtRoutingTable struct {
	self       hash.KadKey
	buckets    [][]peerInfo
	bucketSize int
}

func NewDhtRoutingTable(self hash.KadKey, bucketSize int) *DhtRoutingTable {
	rt := DhtRoutingTable{
		self:       self,
		buckets:    make([][]peerInfo, 0),
		bucketSize: bucketSize,
	}
	// define bucket 0
	rt.buckets = append(rt.buckets, make([]peerInfo, 0))
	return &rt
}

func (rt *DhtRoutingTable) BucketSize() int {
	return rt.bucketSize
}

func (rt *DhtRoutingTable) BucketIdForKey(kadId hash.KadKey) int {
	bid := hash.CommonPrefixLength(rt.self, kadId)
	if bid >= len(rt.buckets) {
		bid = len(rt.buckets) - 1
	}
	return bid
}

func (rt *DhtRoutingTable) SizeOfBucket(bucketId int) int {
	return len(rt.buckets[bucketId])
}

func (rt *DhtRoutingTable) AddPeer(pi peer.AddrInfo) bool {
	return rt.addPeer(hash.PeerKadID(pi.ID), pi)
}

func (rt *DhtRoutingTable) addPeer(kadId hash.KadKey, pi peer.AddrInfo) bool {

	bid := rt.BucketIdForKey(kadId)

	lastBucketId := len(rt.buckets) - 1

	if rt.alreadyInBucket(kadId, bid) {
		// discard new peer
		return false
	}

	if bid < lastBucketId {
		// new peer doesn't belong in last bucket
		if len(rt.buckets[bid]) >= rt.bucketSize {
			// bucket is full, discard new peer
			return false
		}

		// add new peer to bucket
		rt.buckets[bid] = append(rt.buckets[bid], peerInfo{pi, kadId})
		return true
	}
	if len(rt.buckets[lastBucketId]) < rt.bucketSize {
		// last bucket is not full, add new peer
		rt.buckets[lastBucketId] = append(rt.buckets[lastBucketId], peerInfo{pi, kadId})
		return true
	}
	// last bucket is full, try to split it
	for len(rt.buckets[lastBucketId]) == rt.bucketSize {
		// farBucket contains peers with a CPL matching lastBucketId
		farBucket := make([]peerInfo, 0)
		// closeBucket contains peers with a CPL higher than lastBucketId
		closeBucket := make([]peerInfo, 0)

		for _, p := range rt.buckets[lastBucketId] {
			if hash.CommonPrefixLength(p.kadId, rt.self) == lastBucketId {
				farBucket = append(farBucket, p)
			} else {
				closeBucket = append(closeBucket, p)
			}
		}
		if len(farBucket) == rt.bucketSize &&
			hash.CommonPrefixLength(rt.self, kadId) == lastBucketId {
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

	newBid := rt.BucketIdForKey(kadId)
	// add new peer to appropraite bucket
	rt.buckets[newBid] = append(rt.buckets[newBid], peerInfo{pi, kadId})
	return true
}

func (rt *DhtRoutingTable) alreadyInBucket(kadId hash.KadKey, bucketId int) bool {
	for _, p := range rt.buckets[bucketId] {
		if p.kadId == kadId {
			return true
		}
	}
	return false
}

func (rt *DhtRoutingTable) RemovePeer(kadId hash.KadKey) bool {
	bid := rt.BucketIdForKey(kadId)
	for i, p := range rt.buckets[bid] {
		if p.kadId == kadId {
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
	for _, p := range rt.buckets[bid] {
		if p.kadId == kadId {
			return p.id
		}
	}
	return peer.AddrInfo{}
}

// TODO: not working as expected
// returns min(n, bucketSize) peers from the bucket matching the given key
func (rt *DhtRoutingTable) NearestPeers(kadId hash.KadKey, n int) []peer.AddrInfo {
	bid := rt.BucketIdForKey(kadId)
	peers := make([]peerInfo, len(rt.buckets[bid]))
	copy(peers, rt.buckets[bid])
	sort.SliceStable(peers, func(i, j int) bool {
		for k := 0; k < hash.Keysize; k++ {
			distI := peers[i].kadId[k] ^ kadId[k]
			distJ := peers[j].kadId[k] ^ kadId[k]
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
