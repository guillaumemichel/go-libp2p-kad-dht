package network

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO: Use sync.Pool to reuse buffers https://pkg.go.dev/sync#Pool

type MessageEndpoint struct {
	Host host.Host

	writers sync.Pool
	readers sync.Pool
}

func NewMessageEndpoint(host host.Host) *MessageEndpoint {
	return &MessageEndpoint{
		Host:    host,
		writers: sync.Pool{},
		readers: sync.Pool{},
	}
}

func (msgEndpoint *MessageEndpoint) SendDhtRequest(ctx context.Context, p peer.ID, req *pb.DhtMessage) (*pb.DhtMessage, error) {
	s, err := msgEndpoint.Host.NewStream(ctx, p, protocol.ProtocolDHT)
	if err != nil {
		fmt.Println("stream creation error")
		return nil, err
	}
	defer s.Close()

	err = WriteMsg(s, req)
	if err != nil {
		fmt.Println("error writing message")
		return nil, err
	}

	resp := &pb.DhtMessage{}
	err = ReadMsg(s, resp)
	if err != nil {
		fmt.Println("error reading message")
		return nil, err
	}

	return resp, nil
}

func WriteMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	w := pbio.NewDelimitedWriter(s)
	return w.WriteMsg(msg)
}

func ReadMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	r := pbio.NewDelimitedReader(s, network.MessageSizeMax)
	return r.ReadMsg(msg)
}
