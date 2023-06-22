package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	tutil "github.com/libp2p/go-libp2p-kad-dht/examples/util"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/util"
)

var (
	protocolID protocol.ID = "/ipfs/kad/1.0.0"
)

func FindPeer(ctx context.Context) {
	ctx, span := util.StartSpan(ctx, "FindPeer Test")
	defer span.End()

	clk := clock.New()
	sched := simplescheduler.NewSimpleScheduler(clk)

	h, err := tutil.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	pid := peerid.PeerID{ID: h.ID()}
	kadid := pid.Key()

	rt := simplert.NewSimpleRT(kadid, 20)
	msgEndpoint := libp2pendpoint.NewMessageEndpoint(ctx, h, sched)

	friend, err := peer.Decode("12D3KooWGjgvfDkpuVAoNhd7PRRvMTEG4ZgzHBFURqDe1mqEzAMS")
	if err != nil {
		panic(err)
	}
	friendID := peerid.PeerID{ID: friend}

	a, err := multiaddr.NewMultiaddr("/ip4/45.32.75.236/udp/4001/quic")
	if err != nil {
		panic(err)
	}
	friendAddr := peer.AddrInfo{ID: friend, Addrs: []multiaddr.Multiaddr{a}}
	if err := h.Connect(ctx, friendAddr); err != nil {
		panic(err)
	}
	fmt.Println("connected to friend")

	target, err := peer.Decode("12D3KooWMBvV4cphtBLbHQysG6c5nP265aEnXZarCPAHB2UPSGiT")
	if err != nil {
		panic(err)
	}
	targetID := peerid.NewPeerID(target)

	req := ipfskadv1.FindPeerRequest(targetID)
	var resp message.ProtoKadResponseMessage = &ipfskadv1.Message{}
	success, err := rt.AddPeer(ctx, &friendID)
	if err != nil || !success {
		panic("failed to add friend to rt")
	}

	endCond := false
	handleResultsFn := func(ctx context.Context, state simplequery.QueryState,
		_ address.NodeID, resp message.MinKadResponseMessage) simplequery.QueryState {
		msg, ok := resp.(*ipfskadv1.Message)
		if !ok {
			fmt.Println("invalid response!")
			return nil
		}
		fmt.Println("got response! and", len(msg.CloserPeers), "closer peers")
		for _, p := range msg.CloserPeers {
			pid := peer.ID("")
			if pid.UnmarshalBinary(p.Id) != nil {
				fmt.Println("invalid peer id format")
				return nil
			}
			if pid == target {
				fmt.Println("found target!")
				endCond = true
				return nil
			}
		}
		return nil
	}

	simplequery.NewSimpleQuery(ctx, targetID.Key(), address.ProtocolID(protocolID), req, resp, 1, 5*time.Second,
		msgEndpoint, rt, sched, handleResultsFn)

	for i := 0; i < 1000 && !endCond; i++ {
		for sched.RunOne(ctx) {
		}
		time.Sleep(10 * time.Millisecond)
	}
}
