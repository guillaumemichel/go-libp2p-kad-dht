package simplert

import (
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/stretchr/testify/require"
)

func TestBucketId(t *testing.T) {
	zeroBytes := func(n int) []byte {
		bytes := make([]byte, n)
		for i := 0; i < n; i++ {
			bytes[i] = 0
		}
		return bytes
	}
	key0 := hash.KadKey(zeroBytes(32))                          // 00000...000
	key1 := hash.KadKey(append(zeroBytes(31), byte(1)))         // 00000...001
	key2 := hash.KadKey(append([]byte{0x80}, zeroBytes(31)...)) // 10000...000
	key3 := hash.KadKey(append([]byte{0x40}, zeroBytes(31)...)) // 01000...000

	table := NewDhtRoutingTable(key0)
	require.Equal(t, 256, table.BucketIdForKey(key0))
	require.Equal(t, 255, table.BucketIdForKey(key1))
	require.Equal(t, 0, table.BucketIdForKey(key2))
	require.Equal(t, 1, table.BucketIdForKey(key3))
}
