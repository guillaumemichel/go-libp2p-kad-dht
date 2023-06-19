package main

import (
	"context"
	"encoding/hex"
	fmt "fmt"

	"github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	protocolID protocol.ID = "/ipfs/kad/1.0.0"
)

func main() {
	ctx := context.Background()
	h, err := util.Libp2pHost(ctx, "4444")
	if err != nil {
		panic(err)
	}

	friendid, err := peer.Decode("12D3KooWMD2hH3Vr4CNZeTCvkBMGxfeLzdPAXs9mds8z9sMgBxAa")
	if err != nil {
		panic(err)
	}
	a, err := multiaddr.NewMultiaddr("/ip4/94.76.228.174/udp/4001/quic")
	if err != nil {
		panic(err)
	}
	friend := peer.AddrInfo{ID: friendid, Addrs: []multiaddr.Multiaddr{a}}
	if err := h.Connect(ctx, friend); err != nil {
		panic(err)
	}
	fmt.Println("connected to friend")

	marshalled, err := hex.DecodeString("0804122212208c4ebf251d55f0d378c5f2d4ebf5fbfde0c505602c4096b9a7be5200743698e75001")
	if err != nil {
		panic(err)
	}

	req := &Message{}
	err = req.Unmarshal(marshalled)
	if err != nil {
		panic(err)
	}

	fmt.Println(req)

	//ms := NewMessageSenderImpl(h, []protocol.ID{protocolID})
	//resp, err := ms.SendRequest(ctx, friendid, req)

	s, err := h.NewStream(ctx, friendid, protocolID)
	if err != nil {
		panic(err)
	}
	defer s.Close()
}

func WriteMsg1(s network.Stream, msg protoreflect.ProtoMessage) error {
	w := pbio.NewDelimitedWriter(s)
	return w.WriteMsg(msg)
}

func ReadMsg1(s network.Stream, msg protoreflect.ProtoMessage) error {
	r := pbio.NewDelimitedReader(s, network.MessageSizeMax)
	return r.ReadMsg(msg)
}
