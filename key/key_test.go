package key

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	mhreg "github.com/multiformats/go-multihash/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keysize = 32
)

func zeroBytes(n int) []byte {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = 0
	}
	return bytes
}

func randomBytes(n int) []byte {
	blk := make([]byte, n)
	rand.Read(blk)
	return blk
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

func TestKadKeyString(t *testing.T) {
	zeroKadid := KadKey(zeroBytes(keysize))
	zeroHex := strings.Repeat("00", keysize)
	require.Equal(t, zeroHex, zeroKadid.String())

	ffKadid := make([]byte, keysize)
	for i := 0; i < keysize; i++ {
		ffKadid[i] = 0xff
	}
	ffHex := strings.Repeat("ff", keysize)
	require.Equal(t, ffHex, KadKey(ffKadid).String())

	e3Kadid := make([]byte, keysize)
	for i := 0; i < keysize; i++ {
		e3Kadid[i] = 0xe3
	}
	e3Hex := strings.Repeat("e3", keysize)
	require.Equal(t, e3Hex, KadKey(e3Kadid).String())
}

func TestXor(t *testing.T) {
	key0 := KadKey(zeroBytes(keysize))      // 00000...000
	randKey := KadKey(randomBytes(keysize)) // random key

	require.Equal(t, key0, Xor(key0, key0))
	require.Equal(t, randKey, Xor(randKey, key0))
	require.Equal(t, randKey, Xor(key0, randKey))
	require.Equal(t, key0, Xor(randKey, randKey))
}

func TestCommonPrefixLength(t *testing.T) {
	key0 := KadKey(zeroBytes(keysize))                            // 00000...000
	key1 := KadKey(append(zeroBytes(keysize-1), 0x01))            // 00000...001
	key2 := KadKey(append([]byte{0x80}, zeroBytes(keysize-1)...)) // 10000...000
	key3 := KadKey(append([]byte{0x40}, zeroBytes(keysize-1)...)) // 01000...000

	require.Equal(t, keysize*8, CommonPrefixLength(key0, key0))
	require.Equal(t, keysize*8-1, CommonPrefixLength(key0, key1))
	require.Equal(t, 0, CommonPrefixLength(key0, key2))
	require.Equal(t, 1, CommonPrefixLength(key0, key3))
}

func TestCompare(t *testing.T) {
	nKeys := 5
	keys := make([]KadKey, nKeys)
	// ascending order
	keys[0] = KadKey(zeroBytes(keysize))                            // 00000...000
	keys[1] = KadKey(append(zeroBytes(keysize-1), 0x01))            // 00000...001
	keys[2] = KadKey(append(zeroBytes(keysize-1), 0x02))            // 00000...010
	keys[3] = KadKey(append([]byte{0x40}, zeroBytes(keysize-1)...)) // 01000...000
	keys[4] = KadKey(append([]byte{0x80}, zeroBytes(keysize-1)...)) // 10000...000

	for i := 0; i < nKeys; i++ {
		for j := 0; j < nKeys; j++ {
			if i < j {
				require.Equal(t, -1, Compare(keys[i], keys[j]))
			} else if i > j {
				require.Equal(t, 1, Compare(keys[i], keys[j]))
			} else {
				require.Equal(t, 0, Compare(keys[i], keys[j]))
			}
		}
	}
}
