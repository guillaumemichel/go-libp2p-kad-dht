package simplert

import (
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

func TestAdd(t *testing.T) {
	key0 := hash.KadKey(zeroBytes(32))                          // 000000...000
	key1 := hash.KadKey(append([]byte{0x40}, zeroBytes(31)...)) // 010000...000
	key2 := hash.KadKey(append([]byte{0x80}, zeroBytes(31)...)) // 100000...000
	key3 := hash.KadKey(append([]byte{0xc0}, zeroBytes(31)...)) // 110000...000
	key4 := hash.KadKey(append([]byte{0xe0}, zeroBytes(31)...)) // 111000...000
	key5 := hash.KadKey(append([]byte{0x60}, zeroBytes(31)...)) // 011000...000
	key6 := hash.KadKey(append([]byte{0x70}, zeroBytes(31)...)) // 011100...000
	key7 := hash.KadKey(append([]byte{0x18}, zeroBytes(31)...)) // 000110...000
	key8 := hash.KadKey(append([]byte{0x14}, zeroBytes(31)...)) // 000101...000
	key9 := hash.KadKey(append([]byte{0x10}, zeroBytes(31)...)) // 000100...100

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
}
