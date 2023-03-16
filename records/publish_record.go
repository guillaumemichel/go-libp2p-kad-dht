package records

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/multiformats/go-multihash"
)

type PublishRecord struct {
	ID        multihash.Multihash
	ServerKey hash.KadKey
	EncPeerID []byte
	Timestamp uint32
	Signature []byte
}
