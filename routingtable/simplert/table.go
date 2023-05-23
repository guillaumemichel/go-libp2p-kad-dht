package simplert

import (
	"sort"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type RoutingTable interface {
	AddPeer(peer.AddrInfo) bool
	RemovePeer(key.KadKey) bool
	NearestPeers(key.KadKey, int) []peer.ID
	Find(key.KadKey) peer.ID
}

type peerInfo struct {
	id    peer.ID
	kadId key.KadKey
}

type SimpleRT struct {
	self       key.KadKey
	buckets    [][]peerInfo
	bucketSize int
}

func NewSimpleRT(self key.KadKey, bucketSize int) *SimpleRT {
	rt := SimpleRT{
		self:       self,
		buckets:    make([][]peerInfo, 0),
		bucketSize: bucketSize,
	}
	// define bucket 0
	rt.buckets = append(rt.buckets, make([]peerInfo, 0))
	return &rt
}

func (rt *SimpleRT) BucketSize() int {
	return rt.bucketSize
}

func (rt *SimpleRT) BucketIdForKey(kadId key.KadKey) int {
	bid := key.CommonPrefixLength(rt.self, kadId)
	if bid >= len(rt.buckets) {
		bid = len(rt.buckets) - 1
	}
	return bid
}

func (rt *SimpleRT) SizeOfBucket(bucketId int) int {
	return len(rt.buckets[bucketId])
}

func (rt *SimpleRT) AddPeer(p peer.ID) bool {
	return rt.addPeer(key.PeerKadID(p), p)
}

func (rt *SimpleRT) addPeer(kadId key.KadKey, p peer.ID) bool {
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
		rt.buckets[bid] = append(rt.buckets[bid], peerInfo{p, kadId})
		return true
	}
	if len(rt.buckets[lastBucketId]) < rt.bucketSize {
		// last bucket is not full, add new peer
		rt.buckets[lastBucketId] = append(rt.buckets[lastBucketId], peerInfo{p, kadId})
		return true
	}
	// last bucket is full, try to split it
	for len(rt.buckets[lastBucketId]) == rt.bucketSize {
		// farBucket contains peers with a CPL matching lastBucketId
		farBucket := make([]peerInfo, 0)
		// closeBucket contains peers with a CPL higher than lastBucketId
		closeBucket := make([]peerInfo, 0)

		for _, p := range rt.buckets[lastBucketId] {
			if key.CommonPrefixLength(p.kadId, rt.self) == lastBucketId {
				farBucket = append(farBucket, p)
			} else {
				closeBucket = append(closeBucket, p)
			}
		}
		if len(farBucket) == rt.bucketSize &&
			key.CommonPrefixLength(rt.self, kadId) == lastBucketId {
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
	rt.buckets[newBid] = append(rt.buckets[newBid], peerInfo{p, kadId})
	return true
}

func (rt *SimpleRT) alreadyInBucket(kadId key.KadKey, bucketId int) bool {
	for _, p := range rt.buckets[bucketId] {
		if p.kadId == kadId {
			return true
		}
	}
	return false
}

func (rt *SimpleRT) RemovePeer(kadId key.KadKey) bool {
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

func (rt *SimpleRT) Find(kadId key.KadKey) peer.ID {
	bid := rt.BucketIdForKey(kadId)
	for _, p := range rt.buckets[bid] {
		if p.kadId == kadId {
			return p.id
		}
	}
	return peer.ID("")
}

// TODO: not exactly working as expected
// returns min(n, bucketSize) peers from the bucket matching the given key
func (rt *SimpleRT) NearestPeers(kadId key.KadKey, n int) []peer.ID {
	bid := rt.BucketIdForKey(kadId)
	peers := make([]peerInfo, len(rt.buckets[bid]))
	copy(peers, rt.buckets[bid])
	sort.SliceStable(peers, func(i, j int) bool {
		for k := 0; k < key.Keysize; k++ {
			distI := peers[i].kadId[k] ^ kadId[k]
			distJ := peers[j].kadId[k] ^ kadId[k]
			if distI != distJ {
				return distI < distJ
			}
		}
		return false
	})
	pids := make([]peer.ID, min(n, len(peers)))
	for i := 0; i < min(n, len(peers)); i++ {
		pids[i] = peers[i].id
	}

	return pids
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
