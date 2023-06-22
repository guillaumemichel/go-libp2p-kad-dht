package endpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p/core/network"
)

type RequestHandlerFn func(context.Context, address.NodeID,
	message.MinKadMessage) (message.MinKadMessage, error)
type ResponseHandlerFn func(context.Context, message.MinKadResponseMessage, error)

type Endpoint interface {
	// MaybeAddToPeerstore adds the given address to the peerstore if it is
	// valid and if it is not already there.
	MaybeAddToPeerstore(context.Context, address.NetworkAddress, time.Duration) error
	// SendRequestHandleResponse sends a request to the given peer and handles
	// the response with the given handler.
	SendRequestHandleResponse(context.Context, address.ProtocolID, address.NodeID,
		message.MinKadMessage, message.MinKadMessage, ResponseHandlerFn)

	// KadKey returns the KadKey of the local node.
	KadKey() key.KadKey
	//Connectedness(address.NodeID) network.Connectedness
	NetworkAddress(address.NodeID) (address.NetworkAddress, error)
}

type ServerEndpoint interface {
	Endpoint
	// AddRequestHandler registers a handler for a given protocol ID.
	AddRequestHandler(address.ProtocolID, RequestHandlerFn)
	// RemoveRequestHandler removes a handler for a given protocol ID.
	RemoveRequestHandler(address.ProtocolID)
}

type NetworkedEndpoint interface {
	Endpoint
	// Connectedness returns the connectedness of the given peer.
	Connectedness(address.NodeID) network.Connectedness
}
