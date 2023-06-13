package endpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p/core/network"
)

type DialReportFn func(context.Context, bool)
type ResponseHandlerFn func(context.Context, message.MinKadResponseMessage)

type Endpoint interface {
	AsyncDialAndReport(context.Context, address.NodeID, DialReportFn)
	DialPeer(context.Context, address.NodeID) error
	MaybeAddToPeerstore(address.NetworkAddress, time.Duration)
	SendRequest(context.Context, address.NodeID, message.MinKadMessage, message.MinKadMessage) error
	SendRequestHandleResponse(context.Context, address.NodeID, message.MinKadMessage, ResponseHandlerFn)

	// Peerstore functions
	KadID() key.KadKey
	Connectedness(address.NodeID) network.Connectedness
	NetworkAddress(address.NodeID) address.NetworkAddress
}
