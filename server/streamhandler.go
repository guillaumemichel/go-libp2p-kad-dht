package server

import (
	"context"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	endpoint "github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	message "github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	net "github.com/libp2p/go-libp2p/core/network"
)

func (s *Server) DefaultStreamHandler(stream net.Stream) {
	req := &message.Message{}
	err := endpoint.ReadMsg(stream, req)
	if err != nil {
		if err == io.EOF {
			return
		}

		fmt.Println("error reading message:", err)
		stream.Reset()
		return
	}
	events.NewEvent(s.ctx, s.em, func(ctx context.Context) {
		HandleRequest(ctx, s, req, stream)
	})

}
