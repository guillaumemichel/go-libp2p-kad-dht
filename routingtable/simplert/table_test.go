package simplert

import (
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func zeroBytes(n int) []byte {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = 0
	}
	return bytes
}

var (
	key0  = hash.KadKey(zeroBytes(32))                          // 000000...000
	key1  = hash.KadKey(append([]byte{0x40}, zeroBytes(31)...)) // 010000...000
	key2  = hash.KadKey(append([]byte{0x80}, zeroBytes(31)...)) // 100000...000
	key3  = hash.KadKey(append([]byte{0xc0}, zeroBytes(31)...)) // 110000...000
	key4  = hash.KadKey(append([]byte{0xe0}, zeroBytes(31)...)) // 111000...000
	key5  = hash.KadKey(append([]byte{0x60}, zeroBytes(31)...)) // 011000...000
	key6  = hash.KadKey(append([]byte{0x70}, zeroBytes(31)...)) // 011100...000
	key7  = hash.KadKey(append([]byte{0x18}, zeroBytes(31)...)) // 000110...000
	key8  = hash.KadKey(append([]byte{0x14}, zeroBytes(31)...)) // 000101...000
	key9  = hash.KadKey(append([]byte{0x10}, zeroBytes(31)...)) // 000100...000
	key10 = hash.KadKey(append([]byte{0x20}, zeroBytes(31)...)) // 001000...000
	key11 = hash.KadKey(append([]byte{0x30}, zeroBytes(31)...)) // 001100...100
)

func TestBucketSize(t *testing.T) {
	bucketSize := 100
	rt := NewDhtRoutingTable(key0, bucketSize)
	require.Equal(t, bucketSize, rt.BucketSize())
}

func TestAddPeer(t *testing.T) {

	dumbInfo := peer.AddrInfo{}

	rt := NewDhtRoutingTable(key0, 2)

	require.Equal(t, 0, rt.SizeOfBucket(0))

	// add peer CPL=1, bucket=0
	require.True(t, rt.addPeer(key1, dumbInfo))
	require.Equal(t, 1, rt.SizeOfBucket(0))

	// cannot add the same peer twice
	require.False(t, rt.addPeer(key1, dumbInfo))
	require.Equal(t, 1, rt.SizeOfBucket(0))

	// add peer CPL=0, bucket=0
	require.True(t, rt.addPeer(key2, dumbInfo))
	require.Equal(t, 2, rt.SizeOfBucket(0))

	// add peer CPL=0, bucket=0. split of bucket0
	// key1 goes to bucket1
	require.True(t, rt.addPeer(key3, dumbInfo))
	require.Equal(t, 2, rt.SizeOfBucket(0))
	require.Equal(t, 1, rt.SizeOfBucket(1))

	// already 2 peers with CPL = 0, so this should fail
	require.False(t, rt.addPeer(key4, dumbInfo))
	// add peer CPL=1, bucket=1
	require.True(t, rt.addPeer(key5, dumbInfo))
	require.Equal(t, 2, rt.SizeOfBucket(1))

	// already 2 peers with CPL = 1, so this should fail
	// even if bucket 1 is the last bucket
	require.False(t, rt.addPeer(key6, dumbInfo))

	// add two peers with CPL = 3, bucket=2
	require.True(t, rt.addPeer(key7, dumbInfo))
	require.True(t, rt.addPeer(key8, dumbInfo))
	// cannot add a third peer with CPL = 3
	require.False(t, rt.addPeer(key9, dumbInfo))

	// add two peers with CPL = 2, bucket=2
	require.True(t, rt.addPeer(key10, dumbInfo))
	require.True(t, rt.addPeer(key11, dumbInfo))

	// remove all peers with CPL = 0
	rt.RemovePeer(key3)
	rt.RemovePeer(key4)
	// a new peer with CPL = 0 can be added
	require.True(t, rt.AddPeer(dumbInfo))
	// cannot add the same peer twice even tough
	// the bucket is not full
	require.False(t, rt.AddPeer(dumbInfo))
}

func TestRemovePeer(t *testing.T) {
	dumbInfo := peer.AddrInfo{}

	rt := NewDhtRoutingTable(key0, 2)
	rt.addPeer(key1, dumbInfo)
	require.False(t, rt.RemovePeer(key2))
	require.True(t, rt.RemovePeer(key1))
}

func TestFindPeer(t *testing.T) {
	dumbInfo := peer.AddrInfo{ID: "QmPeer"}

	rt := NewDhtRoutingTable(key0, 2)
	rt.addPeer(key1, dumbInfo)
	require.Equal(t, dumbInfo.ID, rt.Find(key1).ID)
	require.Equal(t, peer.ID(""), rt.Find(key2).ID)
	require.True(t, rt.RemovePeer(key1))
	require.Equal(t, peer.ID(""), rt.Find(key1).ID)
}

func TestNearestPeers(t *testing.T) {

	dumbInfo := make([]peer.AddrInfo, 0, 12)
	for i := 0; i < 12; i++ {
		dumbInfo = append(dumbInfo, peer.AddrInfo{ID: peer.ID(fmt.Sprintf("QmPeer%d", i))})
	}

	bucketSize := 5

	rt := NewDhtRoutingTable(key0, bucketSize)
	rt.addPeer(key1, dumbInfo[1])
	rt.addPeer(key2, dumbInfo[2])
	rt.addPeer(key3, dumbInfo[3])
	rt.addPeer(key4, dumbInfo[4])
	rt.addPeer(key5, dumbInfo[5])
	rt.addPeer(key6, dumbInfo[6])
	rt.addPeer(key7, dumbInfo[7])
	rt.addPeer(key8, dumbInfo[8])
	rt.addPeer(key9, dumbInfo[9])
	rt.addPeer(key10, dumbInfo[10])
	rt.addPeer(key11, dumbInfo[11])

	// find the 2 nearest peers to key0
	peers := rt.NearestPeers(key0, 10)
	require.Equal(t, bucketSize, len(peers))

	keys := make([]peer.ID, 0, len(peers))
	for _, p := range peers {
		keys = append(keys, p.ID)
	}
	expectedOrder := []peer.ID{dumbInfo[9].ID, dumbInfo[8].ID, dumbInfo[7].ID, dumbInfo[10].ID, dumbInfo[11].ID}
	require.Equal(t, expectedOrder, keys)

	peers = rt.NearestPeers(key11, 2)
	require.Equal(t, 2, len(peers))

	// create routing table with a single duplicate peer
	// useful to test peers sorting with duplicate (even tough it should never happen)
	rt2 := NewDhtRoutingTable(key0, 2)
	rt2.buckets[0] = append(rt2.buckets[0], peerInfo{dumbInfo[1], key1})
	rt2.buckets[0] = append(rt2.buckets[0], peerInfo{dumbInfo[1], key1})
	peers = rt2.NearestPeers(key0, 10)
	require.Equal(t, peers[0], peers[1])
}
