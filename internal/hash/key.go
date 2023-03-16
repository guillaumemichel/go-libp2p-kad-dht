package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multihash"
	mhreg "github.com/multiformats/go-multihash/core"
)

const (
	// HasherID is the identifier hash function used to derive the second hash identifiers
	// associated with a CID or multihash
	HasherID = multihash.SHA2_256

	// Keysize is the length in bytes of the hash function's digest, which is equivalent to the keysize in the Kademlia keyspace
	Keysize = sha256.Size
)

type KadKey [Keysize]byte

func HexKadID(kadid KadKey) string {
	return hex.EncodeToString(kadid[:])
}

func PeerKadID(p peer.ID) KadKey {
	// hasher is the hash function used to derive the second hash identifiers
	hasher, _ := mhreg.GetHasher(HasherID)
	hasher.Write([]byte(p))
	return KadKey(hasher.Sum(nil))
}
