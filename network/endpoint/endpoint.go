package endpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p/core/network"
)

type DialReportFn func(context.Context, bool)
type ResponseHandlerFn func(context.Context, message.MinKadResponseMessage, error)

type Endpoint interface {
	AsyncDialAndReport(context.Context, address.NodeID, DialReportFn)
	DialPeer(context.Context, address.NodeID) error
	MaybeAddToPeerstore(context.Context, address.NetworkAddress, time.Duration) error
	SendRequest(context.Context, address.NodeID, message.MinKadMessage, message.MinKadMessage) error
	SendRequestHandleResponse(context.Context, address.NodeID, message.MinKadMessage, ResponseHandlerFn)

	// Peerstore functions
	KadKey() key.KadKey
	Connectedness(address.NodeID) network.Connectedness
	NetworkAddress(address.NodeID) (address.NetworkAddress, error)
}
