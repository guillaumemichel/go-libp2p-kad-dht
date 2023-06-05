package server

import (
	"context"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	endpoint "github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	net "github.com/libp2p/go-libp2p/core/network"
)

func (s *Server) DefaultStreamHandler(stream net.Stream) {
	req := &pb.Message{}
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
