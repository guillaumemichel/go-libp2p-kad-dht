package server

import (
	"context"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	net "github.com/libp2p/go-libp2p/core/network"
)

func (s *Server) DefaultStreamHandler(stream net.Stream) {
	req := &pb.Message{}
	err := network.ReadMsg(stream, req)
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
