package key

import (
	"crypto/sha256"
	"encoding/hex"
	"math"

	"github.com/libp2p/go-libp2p/core/peer"
	mh "github.com/multiformats/go-multihash"
	mhreg "github.com/multiformats/go-multihash/core"
)

const (
	// HasherID is the identifier hash function used to derive the second hash identifiers
	// associated with a CID or multihash
	HasherID = mh.SHA2_256

	// Keysize is the length in bytes of the hash function's digest, which is equivalent to the keysize in the Kademlia keyspace
	Keysize = sha256.Size
)

type KadKey [Keysize]byte

var (
	ZeroKey = KadKey{}
)

func PeerKadID(p peer.ID) KadKey {
	return StringKadID(string(p))
}

func StringKadID(s string) KadKey {
	// hasher is the hash function used to derive the second hash identifiers
	hasher, _ := mhreg.GetHasher(HasherID)
	hasher.Write([]byte(s))
	return KadKey(hasher.Sum(nil))
}

func (k KadKey) Hex() string {
	return hex.EncodeToString(k[:])
}

func (k KadKey) String() string {
	return k.Hex()
}

func (k KadKey) Xor(other KadKey) KadKey {
	var xored KadKey
	for i := 0; i < Keysize; i++ {
		xored[i] = k[i] ^ other[i]
	}
	return xored
}

func CommonPrefixLength(a, b KadKey) int {
	var xored byte
	for i := 0; i < Keysize; i++ {
		xored = a[i] ^ b[i]
		if xored != 0 {
			return i*8 + 7 - int(math.Log2(float64(xored)))
		}
	}
	return 8 * Keysize

}

func (k KadKey) Compare(other KadKey) int {
	for i := 0; i < Keysize; i++ {
		if k[i] < other[i] {
			return -1
		}
		if k[i] > other[i] {
			return 1
		}
	}
	return 0
}
