package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p/core/peer"
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
	protocolID address.ProtocolID = "/ipfs/kad/1.0.0" // IPFS DHT network protocol ID
)

func FindPeer(ctx context.Context) {
	ctx, span := util.StartSpan(ctx, "FindPeer Test")
	defer span.End()

	// this example is using real time
	clk := clock.New()

	// create a libp2p host
	h, err := tutil.Libp2pHost(ctx, "8888")
	if err != nil {
		panic(err)
	}

	pid := peerid.NewPeerID(h.ID())
	// get the peer's kademlia key (derived from its peer.ID)
	kadid := pid.Key()

	// create a simple routing table, with bucket size 20
	rt := simplert.NewSimpleRT(kadid, 20)
	// create a scheduler using real time
	sched := simplescheduler.NewSimpleScheduler(clk)
	// create a message endpoint is used to communicate with other peers
	msgEndpoint := libp2pendpoint.NewMessageEndpoint(ctx, h, sched)

	// friend is the first peer we know in the IPFS DHT network (bootstrap node)
	friend, err := peer.Decode("12D3KooWGjgvfDkpuVAoNhd7PRRvMTEG4ZgzHBFURqDe1mqEzAMS")
	if err != nil {
		panic(err)
	}
	friendID := peerid.NewPeerID(friend)

	// multiaddress of friend
	a, err := multiaddr.NewMultiaddr("/ip4/45.32.75.236/udp/4001/quic")
	if err != nil {
		panic(err)
	}
	// connect to friend
	friendAddr := peer.AddrInfo{ID: friend, Addrs: []multiaddr.Multiaddr{a}}
	if err := h.Connect(ctx, friendAddr); err != nil {
		panic(err)
	}
	fmt.Println("connected to friend")

	// target is the peer we want to find
	target, err := peer.Decode("12D3KooWMBvV4cphtBLbHQysG6c5nP265aEnXZarCPAHB2UPSGiT")
	if err != nil {
		panic(err)
	}
	targetID := peerid.NewPeerID(target)

	// create a find peer request message
	req := ipfskadv1.FindPeerRequest(targetID)
	// empty response message to be filled by the query process, the protobuf
	// message must be know to parse the response
	var resp message.ProtoKadResponseMessage = &ipfskadv1.Message{}
	// add friend to routing table
	success, err := rt.AddPeer(ctx, friendID)
	if err != nil || !success {
		panic("failed to add friend to rt")
	}

	// endCond is used to terminate the simulation once the query is done
	endCond := false
	handleResultsFn := func(ctx context.Context, state simplequery.QueryState,
		id address.NodeID, resp message.MinKadResponseMessage) (bool, simplequery.QueryState) {
		// parse response to ipfs dht message
		msg, ok := resp.(*ipfskadv1.Message)
		if !ok {
			fmt.Println("invalid response!")
			return false, nil
		}
		peers := make([]address.NodeID, 0, len(msg.CloserPeers))
		for _, p := range msg.CloserPeers {
			pid := peer.ID("")
			if pid.UnmarshalBinary(p.Id) != nil {
				fmt.Println("invalid peer id format")
				continue
			}
			peers = append(peers, peerid.NewPeerID(pid))
			if pid == target {
				endCond = true
			}
		}
		fmt.Println("---\nResponse from", id, "with", peers)
		if endCond {
			fmt.Println("  - target found!", target)
		}
		return endCond, nil
	}

	// create the query, using the target kademlia key as target, the IPFS DHT
	// protocol ID, the request and response messages in the IPFS DHT format,
	// a concurrency parameter of 1, a timeout of 5 seconds, the libp2p message
	// endpoint, the node's routing table and scheduler, and the response
	// handler function.
	// The query will be executed only once actions are run on the scheduler.
	// For now, it is only scheduled to be run.
	simplequery.NewSimpleQuery(ctx, targetID.Key(), protocolID, req, resp, 1,
		5*time.Second, msgEndpoint, rt, sched, handleResultsFn)

	span.AddEvent("start request execution")

	// run the actions from the scheduler until the query is done
	for i := 0; i < 1000 && !endCond; i++ {
		for sched.RunOne(ctx) {
		}
		time.Sleep(10 * time.Millisecond)
	}
}
