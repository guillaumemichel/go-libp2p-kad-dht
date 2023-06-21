package simserver

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
)

type ReplyFn func(message.MinKadResponseMessage)

type SimServer interface {
	HandleFindNodeRequest(context.Context, address.NetworkAddress,
		message.MinKadMessage, ReplyFn)
}
