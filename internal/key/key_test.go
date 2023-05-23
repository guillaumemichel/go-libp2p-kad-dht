package key

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	mhreg "github.com/multiformats/go-multihash/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func zeroBytes(n int) []byte {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = 0
	}
	return bytes
}

func TestPeerKadID(t *testing.T) {
	peerid := peer.ID("12D3KooWFHYCmTCexEziKBVFEVT4FAVPBUYJKEBUrdgu9RYoVE1T")
	kadid := PeerKadID(peerid)

	// get sha256 hasher
	hasher, err := mhreg.GetHasher(HasherID)
	assert.NoError(t, err)
	hasher.Write([]byte(peerid))
	expectedKadid := hasher.Sum(nil)

	assert.Equal(t, kadid[:], expectedKadid)
}

func TestCommonPrefixLength(t *testing.T) {
	key0 := KadKey(zeroBytes(32))                          // 00000...000
	key1 := KadKey(append(zeroBytes(31), byte(1)))         // 00000...001
	key2 := KadKey(append([]byte{0x80}, zeroBytes(31)...)) // 10000...000
	key3 := KadKey(append([]byte{0x40}, zeroBytes(31)...)) // 01000...000

	require.Equal(t, 256, CommonPrefixLength(key0, key0))
	require.Equal(t, 255, CommonPrefixLength(key0, key1))
	require.Equal(t, 0, CommonPrefixLength(key0, key2))
	require.Equal(t, 1, CommonPrefixLength(key0, key3))
}
