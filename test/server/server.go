package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server"
	"github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

func main() {
	pi := serv()
	fmt.Println(pi)
	client(pi)
}

func serv() peer.AddrInfo {
	ctx := context.Background()
	h, err := util.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	em := events.NewEventsManager()
	rt := simplert.NewSimpleRT(key.PeerKadID(h.ID()), 20)
	serv := server.NewServer(h, rt, em)
	server.SetStreamHandler(serv, serv.DefaultStreamHandler, consts.ProtocolDHT)

	p := peer.ID("12D3KooWG2qAjJvJwv4K7hrHbNVJdDzQqqwPSEezM1R3csV22yK3")
	h.Peerstore().AddAddr(p, multiaddr.StringCast("/ip4/1.2.3.4/tcp/5678"), peerstore.PermanentAddrTTL)
	rt.AddPeer(p)

	return h.Peerstore().PeerInfo(h.ID())
}

func client(serv peer.AddrInfo) {
	ctx := context.Background()
	h, err := util.Libp2pHost(ctx, "9999")
	if err != nil {
		panic(err)
	}

	p := peer.ID("12D3KooWG2qAjJvJwv4K7hrHbNVJdDzQqqwPSEezM1R3csV22yK3")

	msg := &pb.Message{
		Type: pb.Message_FIND_NODE,
		Key:  []byte(p),
	}

	serv = peer.AddrInfo{
		ID:    serv.ID,
		Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/0.0.0.0/tcp/8888")},
	}
	err = h.Connect(ctx, serv)
	if err != nil {
		panic(err)
	}
	stream, err := h.NewStream(ctx, serv.ID, consts.ProtocolDHT)
	if err != nil {
		panic(err)
	}

	err = network.WriteMsg(stream, msg)
	if err != nil {
		panic(err)
	}
	msg = &pb.Message{}
	err = network.ReadMsg(stream, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("we got a response!")
	peers := msg.GetCloserPeers()
	peerid := string(peers[0].Id)
	var maddr multiaddr.Multiaddr
	for _, a := range peers[0].Addrs {
		maddr, err = multiaddr.NewMultiaddrBytes(a)
		if err != nil {
			panic(err)
		}
		break
	}
	fmt.Println(peerid, maddr)
}
