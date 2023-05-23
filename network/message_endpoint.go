package network

import (
	"context"
	"fmt"
	"sync"
	"time"

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

// TODO: totally change this function
// Timeout should be handled by this function, maybe with context??
func (msgEndpoint *MessageEndpoint) SendDhtRequest(p peer.ID, req *pb.Message, timeout time.Duration) error {
	ctx := context.Background()                                           // TODO: figure out context
	s, err := msgEndpoint.Host.NewStream(ctx, p, "/dummy/protocol/1.0.0") // TODO: update protocol
	if err != nil {
		fmt.Println("stream creation error")
		return err
	}
	defer s.Close()

	err = WriteMsg(s, req)
	if err != nil {
		fmt.Println("error writing message")
		return err
	}

	resp := &pb.Message{}
	err = ReadMsg(s, resp)
	if err != nil {
		fmt.Println("error reading message")
		return err
	}

	return nil
}

func WriteMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	w := pbio.NewDelimitedWriter(s)
	return w.WriteMsg(msg)
}

func ReadMsg(s network.Stream, msg protoreflect.ProtoMessage) error {
	r := pbio.NewDelimitedReader(s, network.MessageSizeMax)
	return r.ReadMsg(msg)
}
