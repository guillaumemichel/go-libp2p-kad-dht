package testserver

import (
	"context"
	"fmt"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	endpoint "github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	message "github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server"
	tutil "github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/libp2p/go-libp2p-kad-dht/util"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multibase"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

var (
	targetBytesID = "mACQIARIgp9PBu+JuU8aicuW8xT+Oa08OntMyqdLbfQtOplAHlME"
)

func ServerTest(ctx context.Context) {
	newCtx, span := util.StartSpan(ctx, "ServerTest")
	clk := clock.New()
	ai := serv(newCtx, clk)
	client(newCtx, ai, clk)
	span.End()
}

func serv(ctx context.Context, clk clock.Clock) peer.AddrInfo {
	h, err := tutil.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	sched := simplescheduler.NewSimpleScheduler(ctx, clk)
	rt := simplert.NewSimpleRT(key.PeerKadID(h.ID()), 20)
	serv := server.NewServer(ctx, h, rt, sched, []address.ProtocolID{consts.ProtocolDHT})
	server.SetStreamHandler(serv, serv.DefaultStreamHandler, consts.ProtocolDHT)

	//p := peer.ID("12D3KooWG2qAjJvJwv4K7hrHbNVJdDzQqqwPSEezM1R3csV22yK3")
	_, bin, _ := multibase.Decode(targetBytesID)
	p := peer.ID(bin)
	h.Peerstore().AddAddr(p, multiaddr.StringCast("/ip4/1.2.3.4/tcp/5678"), peerstore.PermanentAddrTTL)
	rt.AddPeer(ctx, p)

	return h.Peerstore().PeerInfo(h.ID())
}

func client(ctx context.Context, serv peer.AddrInfo, clk clock.Clock) {
	h, err := tutil.Libp2pHost(ctx, "9999")
	if err != nil {
		panic(err)
	}
	_, bin, _ := multibase.Decode(targetBytesID)
	p := peer.ID(bin)
	marshalledPeerid, _ := p.MarshalBinary()
	msg := &message.Message{
		Type: message.Message_FIND_NODE,
		Key:  marshalledPeerid,
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

	err = endpoint.WriteMsg(stream, msg)
	if err != nil {
		panic(err)
	}
	msg = &message.Message{}
	err = endpoint.ReadMsg(stream, msg)
	if err != nil {
		panic(err)
	}

	peers := msg.GetCloserPeers()

	peerid := peer.ID("")
	err = peerid.UnmarshalBinary(peers[0].Id)
	if err != nil {
		panic(err)
	}

	maddr, err := multiaddr.NewMultiaddrBytes(peers[0].Addrs[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(peerid, maddr)
}
