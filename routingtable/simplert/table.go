package simplert

import (
	"context"
	"sort"
	"strconv"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type peerInfo struct {
	id    address.NodeID
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

func (rt *SimpleRT) AddPeer(ctx context.Context, id address.NodeID) bool {
	return rt.addPeer(ctx, address.KadID(id), id)
}

func (rt *SimpleRT) addPeer(ctx context.Context, kadId key.KadKey, id address.NodeID) bool {
	_, span := internal.StartSpan(ctx, "simplert.addPeer", trace.WithAttributes(
		attribute.String("KadID", kadId.Hex()),
		attribute.Stringer("PeerID", id),
	))
	defer span.End()

	bid := rt.BucketIdForKey(kadId)

	lastBucketId := len(rt.buckets) - 1

	if rt.alreadyInBucket(kadId, bid) {
		span.AddEvent("peer not added, already in bucket " + strconv.Itoa(bid))
		// discard new peer
		return false
	}

	if bid < lastBucketId {
		// new peer doesn't belong in last bucket
		if len(rt.buckets[bid]) >= rt.bucketSize {
			span.AddEvent("peer not added, bucket " + strconv.Itoa(bid) + " full")
			// bucket is full, discard new peer
			return false
		}

		// add new peer to bucket
		rt.buckets[bid] = append(rt.buckets[bid], peerInfo{id, kadId})
		span.AddEvent("peer added to bucket " + strconv.Itoa(bid))
		return true
	}
	if len(rt.buckets[lastBucketId]) < rt.bucketSize {
		// last bucket is not full, add new peer
		rt.buckets[lastBucketId] = append(rt.buckets[lastBucketId], peerInfo{id, kadId})
		span.AddEvent("peer added to bucket " + strconv.Itoa(lastBucketId))
		return true
	}
	// last bucket is full, try to split it
	for len(rt.buckets[lastBucketId]) == rt.bucketSize {
		// farBucket contains peers with a CPL matching lastBucketId
		farBucket := make([]peerInfo, 0)
		// closeBucket contains peers with a CPL higher than lastBucketId
		closeBucket := make([]peerInfo, 0)

		span.AddEvent("splitting last bucket (" + strconv.Itoa(lastBucketId) + ")")

		for _, p := range rt.buckets[lastBucketId] {
			if key.CommonPrefixLength(p.kadId, rt.self) == lastBucketId {
				farBucket = append(farBucket, p)
			} else {
				closeBucket = append(closeBucket, p)
				span.AddEvent(p.id.String() + " moved to new bucket (" +
					strconv.Itoa(lastBucketId+1) + ")")
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
	rt.buckets[newBid] = append(rt.buckets[newBid], peerInfo{id, kadId})
	span.AddEvent("peer added to bucket " + strconv.Itoa(newBid))
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

func (rt *SimpleRT) RemovePeer(ctx context.Context, kadId key.KadKey) bool {
	_, span := internal.StartSpan(ctx, "simplert.removePeer", trace.WithAttributes(
		attribute.String("KadID", kadId.Hex()),
	))
	defer span.End()

	bid := rt.BucketIdForKey(kadId)
	for i, p := range rt.buckets[bid] {
		if p.kadId == kadId {
			// remove peer from bucket
			rt.buckets[bid][i] = rt.buckets[bid][len(rt.buckets[bid])-1]
			rt.buckets[bid] = rt.buckets[bid][:len(rt.buckets[bid])-1]

			span.AddEvent(p.id.String() + " removed from bucket " + strconv.Itoa(bid))
			return true
		}
	}
	// peer not found in the routing table
	span.AddEvent("peer not found in bucket " + strconv.Itoa(bid))
	return false
}

func (rt *SimpleRT) Find(kadId key.KadKey) address.NodeID {
	bid := rt.BucketIdForKey(kadId)
	for _, p := range rt.buckets[bid] {
		if p.kadId == kadId {
			return p.id
		}
	}
	return nil
}

// TODO: not exactly working as expected
// returns min(n, bucketSize) peers from the bucket matching the given key
func (rt *SimpleRT) NearestPeers(ctx context.Context, kadId key.KadKey, n int) []address.NodeID {
	_, span := internal.StartSpan(ctx, "simplert.nearestPeers", trace.WithAttributes(
		attribute.String("KadID", kadId.Hex()),
		attribute.Int("n", n),
	))
	defer span.End()

	bid := rt.BucketIdForKey(kadId)

	var peers []peerInfo
	// TODO: optimize this
	if len(rt.buckets[bid]) == n {
		peers = make([]peerInfo, len(rt.buckets[bid]))
		copy(peers, rt.buckets[bid])
	} else {
		peers = make([]peerInfo, 0)
		for i := 0; i < len(rt.buckets); i++ {
			for _, p := range rt.buckets[i] {
				if p.kadId != rt.self {
					peers = append(peers, p)
				}
			}
		}
	}

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
	pids := make([]address.NodeID, min(n, len(peers)))
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
