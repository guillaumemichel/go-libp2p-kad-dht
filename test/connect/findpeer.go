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
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/libp2pendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	tutil "github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/libp2p/go-libp2p-kad-dht/util"
)

var (
	protocolID protocol.ID = "/ipfs/kad/1.0.0"
)

func FindPeer(ctx context.Context) {
	ctx, span := util.StartSpan(ctx, "FindPeer Test")
	defer span.End()

	clk := clock.New()
	sched := simplescheduler.NewSimpleScheduler(ctx, clk)

	h, err := tutil.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	peerid := h.ID()
	kadid := key.PeerKadID(peerid)

	rt := simplert.NewSimpleRT(kadid, 20)
	msgEndpoint := libp2pendpoint.NewMessageEndpoint(h, protocolID)

	friendid, err := peer.Decode("12D3KooWGjgvfDkpuVAoNhd7PRRvMTEG4ZgzHBFURqDe1mqEzAMS")
	if err != nil {
		panic(err)
	}
	a, err := multiaddr.NewMultiaddr("/ip4/45.32.75.236/udp/4001/quic")
	if err != nil {
		panic(err)
	}
	friend := peer.AddrInfo{ID: friendid, Addrs: []multiaddr.Multiaddr{a}}
	if err := h.Connect(ctx, friend); err != nil {
		panic(err)
	}
	fmt.Println("connected to friend")

	target, err := peer.Decode("QmXnMRxi3V4W8j2ZCHem2QSJoiYm7wRFzzVyZ6rZeCVmqg")
	if err != nil {
		panic(err)
	}

	req := ipfskadv1.FindPeerRequest(target)
	if !rt.AddPeer(ctx, friendid) {
		panic("failed to add friend to rt")
	}
	//req.ClusterLevelRaw = int32(1)

	endCond := false
	handleResultsFn := func(ctx context.Context, state simplequery.QueryState, resp message.MinKadResponseMessage) simplequery.QueryState {
		msg, ok := resp.(*ipfskadv1.Message)
		if !ok {
			fmt.Println("invalid response!")
			return nil
		}
		fmt.Println("number of peers returned", len(msg.CloserPeers))
		for _, p := range msg.CloserPeers {
			pid := peer.ID("")
			if pid.UnmarshalBinary(p.Id) != nil {
				fmt.Println("invalid peer id format")
				return nil
			}
			fmt.Println(pid)
			if pid == target {
				fmt.Println("found target!")
				endCond = true
				return nil
			}
		}
		return nil
	}

	simplequery.NewSimpleQuery(ctx, key.PeerKadID(friendid), req, 1, 5*time.Second, msgEndpoint, rt, sched, handleResultsFn)

	for !endCond {
		for sched.RunOne(ctx) {
		}
		time.Sleep(10 * time.Millisecond)
	}
}
