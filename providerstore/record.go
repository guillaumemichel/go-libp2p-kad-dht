package providerstore

import (
	"github.com/libp2p/go-libp2p-kad-dht/internal/hashing"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ProviderRecord struct {
	ServerKey hashing.KadKey
	Provider  peer.ID
	EncPeerID []byte
	Timestamp uint32 // TODO: define a proper timestamp structure
	Signature []byte
}
