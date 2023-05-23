package consts

import "github.com/libp2p/go-libp2p/core/protocol"

const (
	// BucketSize is the default bucket size for the DHT.
	BucketSize = 20
	// NClosestPeers is the default number of closest peers returned by the DHT.
	NClosestPeers = 20
	// ReplicationFactor is the default replication factor for the DHT.
	ReplicationFactor = 20
)

var (
	// ProtocolDHT is the default DHT protocol.
	ProtocolDHT protocol.ID = "/ipfs/kad/1.1.0"
	// DefaultProtocols spoken by the DHT.
	DefaultProtocols = []protocol.ID{ProtocolDHT}
)
