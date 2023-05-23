package server

import (
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	net "github.com/libp2p/go-libp2p/core/network"
)

func (s *Server) DefaultStreamHandler(stream net.Stream) {
	for {
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

		// TODO: improve this
		event := struct {
			*pb.Message
			net.Stream
		}{req, stream}

		events.NewEvent(s.em, event)

	}
}
