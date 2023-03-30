package routing

import (
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
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
	keys = []hash.KadKey{
		hash.KadKey(zeroBytes(hash.Keysize)),                // 0000 0000 ... 0000 0000
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x1)), // 0000 0000 ... 0000 0001
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x2)), // 0000 0000 ... 0000 0010
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x3)), // 0000 0000 ... 0000 0011
	}
)

func TestQpeersetInsert(t *testing.T) {
	q := newQpeerset(keys[0])
	p0 := &qpeer{dist: keys[0]}
	p1 := &qpeer{dist: keys[1]}
	p2 := &qpeer{dist: keys[2]}
	p3 := &qpeer{dist: keys[3]}

	require.Equal(t, 0, q.size)
	q.insert(p0)
	require.Equal(t, 1, q.size)
	q.insert(p0)
	require.Equal(t, 1, q.size)
	q.insert(p1)
	require.Equal(t, 2, q.size)
	require.Equal(t, p0.dist, q.head.dist)
	require.Equal(t, p1.dist, q.head.next.dist)

	q.insert(p3)
	q.insert(p2)
	require.Equal(t, 4, q.size)
	require.Equal(t, p2.dist, q.head.next.next.dist)
}
