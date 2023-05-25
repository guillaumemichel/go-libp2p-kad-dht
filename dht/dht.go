package dht

import (
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type KademliaDHT struct {
	rt routingtable.RoutingTable

	protocols []protocol.ID

	serverProtocols []protocol.ID
}
