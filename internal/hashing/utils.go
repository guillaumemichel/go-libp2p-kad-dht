package hashing

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/multiformats/go-multihash"
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
