package server

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
)

type Server interface {
	HandleRequest(context.Context, address.NodeID,
		message.MinKadMessage) (message.MinKadMessage, error)
}

//var _ endpoint.RequestHandlerFn = (Server)(nil).HandleRequest
