package endpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type DialReportFn func(context.Context, bool)

type Endpoint interface {
	AsyncDialAndReport(context.Context, address.NodeID, DialReportFn)
	DialPeer(context.Context, address.NodeID) error
	MaybeAddToPeerstore(address.NetworkAddress, time.Duration)
	SendRequest(context.Context, address.NodeID, message.MinKadRequestMessage, message.MinKadResponseMessage, protocol.ID) error

	// Peerstore functions
	KadID() key.KadKey
	Connectedness(address.NodeID) network.Connectedness
	PeerInfo(address.NodeID) peer.AddrInfo
}
