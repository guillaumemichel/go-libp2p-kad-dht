package consts

import (
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
)

const (
	// BucketSize is the default bucket size for the DHT.
	BucketSize = 20
	// NClosestPeers is the default number of closest peers returned by the DHT.
	NClosestPeers = 20
	// ReplicationFactor is the default replication factor for the DHT.
	ReplicationFactor = 20

	// PeerstoreTTL is the default TTL for an entry in the peerstore.
	PeerstoreTTL = 10 * time.Minute
)

var (
	// ProtocolDHT is the default DHT protocol.
	ProtocolDHT address.ProtocolID = "/ipfs/kad/1.1.0"
	// DefaultProtocols spoken by the DHT.
	DefaultProtocols = []address.ProtocolID{ProtocolDHT}
)
