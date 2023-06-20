package main

import (
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/kadid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/fakeendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	sq "github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver/kadsimserver"
	"github.com/libp2p/go-libp2p-kad-dht/util"

	sd "github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	ss "github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
)

const (
	keysize      = 32
	peerstoreTTL = 10 * time.Minute
)

func zeroBytes(n int) []byte {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = 0
	}
	return bytes
}

func findNode(ctx context.Context) {
	ctx, span := util.StartSpan(ctx, "findNode test")
	defer span.End()

	clk := clock.NewMock()

	dispatcher := sd.NewSimpleDispatcher(clk)

	nodeCount := 4
	ids := make([]kadid.KadID, nodeCount)
	ids[0] = kadid.KadID{KadKey: key.KadKey(zeroBytes(keysize))}
	ids[1] = kadid.KadID{KadKey: key.KadKey(append(zeroBytes(keysize-1), 0x01))}
	ids[2] = kadid.KadID{KadKey: key.KadKey(append(zeroBytes(keysize-1), 0x02))}
	ids[3] = kadid.KadID{KadKey: key.KadKey(append(zeroBytes(keysize-1), 0x03))}

	//     ^
	//    / \
	//   ^   ^
	//  A B C D

	rts := make([]*simplert.SimpleRT, len(ids))
	eps := make([]*fakeendpoint.FakeEndpoint, len(ids))
	schedulers := make([]*ss.SimpleScheduler, len(ids))
	servers := make([]*kadsimserver.KadSimServer, len(ids))

	for i := 0; i < len(ids); i++ {
		rts[i] = simplert.NewSimpleRT(ids[i].KadKey, 2)
		eps[i] = fakeendpoint.NewFakeEndpoint(ids[i], dispatcher)
		schedulers[i] = ss.NewSimpleScheduler(ctx, clk)
		servers[i] = kadsimserver.NewKadSimServer(rts[i], eps[i])
		dispatcher.AddPeer(ids[i], schedulers[i], servers[i])
	}

	// A connects to B
	eps[0].MaybeAddToPeerstore(ctx, ids[1], peerstoreTTL)
	rts[0].AddPeer(ctx, ids[1])
	eps[1].MaybeAddToPeerstore(ctx, ids[0], peerstoreTTL)
	rts[1].AddPeer(ctx, ids[0])

	// B connects to C
	eps[1].MaybeAddToPeerstore(ctx, ids[2], peerstoreTTL)
	rts[1].AddPeer(ctx, ids[2])
	eps[2].MaybeAddToPeerstore(ctx, ids[1], peerstoreTTL)
	rts[2].AddPeer(ctx, ids[1])

	// C connects to D
	eps[2].MaybeAddToPeerstore(ctx, ids[3], peerstoreTTL)
	rts[2].AddPeer(ctx, ids[3])
	eps[3].MaybeAddToPeerstore(ctx, ids[2], peerstoreTTL)
	rts[3].AddPeer(ctx, ids[2])

	req := simmessage.NewSimRequest(ids[3].Key())
	resp := &simmessage.SimMessage{}

	handleResFn := func(_ context.Context, _ sq.QueryState, _ address.NodeID, msg message.MinKadResponseMessage) sq.QueryState {
		resp, ok := msg.(*simmessage.SimMessage)
		if !ok {
			panic("not a simmessage")
		}
		fmt.Println("got response")
		for _, peer := range resp.CloserNodes() {
			fmt.Println(peer)
			if peer.NodeID().String() == ids[3].NodeID().String() {
				fmt.Println("success", peer)
			}
		}
		return nil
	}

	sq.NewSimpleQuery(ctx, ids[3].Key(), req, resp, 1, time.Second, eps[0], rts[0], schedulers[0], handleResFn)

	dispatcher.DispatchLoop(ctx)

}
