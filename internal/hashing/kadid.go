package hashing

import (
	"github.com/libp2p/go-libp2p/core/peer"
	mhreg "github.com/multiformats/go-multihash/core"
)

func PeerKadID(p peer.ID) KadKey {
	// hasher is the hash function used to derive the second hash identifiers
	hasher, _ := mhreg.GetHasher(HasherID)
	hasher.Write([]byte(p))
	return KadKey(hasher.Sum(nil))
}
