package address

import "github.com/libp2p/go-libp2p-kad-dht/key"

type NetworkAddress interface {
	NodeID() NodeID
}

type NodeID interface {
	NodeID() NodeID
	Key() key.KadKey
	// String returns the string representation of the NodeID. String
	// representation should be unique for each NodeID.
	String() string
}

type ProtocolID string
