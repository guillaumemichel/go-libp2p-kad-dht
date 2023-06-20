package fakeendpoint

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/libp2p/go-libp2p-kad-dht/events/dispatch/simpledispatcher"
	"github.com/libp2p/go-libp2p-kad-dht/events/scheduler/simplescheduler"
	si "github.com/libp2p/go-libp2p-kad-dht/network/address/stringid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server/simserver/kadsimserver"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"
)

var peerstoreTTL = time.Minute

func TestFakeEndpoint(t *testing.T) {
	ctx := context.Background()
	clk := clock.NewMock()

	kadid := si.StringID("self")
	dispatcher := simpledispatcher.NewSimpleDispatcher(clk)
	fakeEndpoint := NewFakeEndpoint(kadid, dispatcher)

	c, err := kadid.Key().Compare(fakeEndpoint.KadKey())
	require.NoError(t, err)
	require.Equal(t, int8(0), c)

	node0 := si.StringID("node0")
	err = fakeEndpoint.DialPeer(ctx, node0)
	require.Equal(t, endpoint.ErrUnknownPeer, err)

	connectedness := fakeEndpoint.Connectedness(node0)
	require.Equal(t, network.NotConnected, connectedness)

	_, err = fakeEndpoint.NetworkAddress(node0)
	require.Equal(t, endpoint.ErrUnknownPeer, err)

	req := simmessage.NewSimRequest(kadid.Key())
	resp := &simmessage.SimMessage{}

	var runCheck bool
	respHandler := func(ctx context.Context, msg message.MinKadResponseMessage, err error) {
		require.NoError(t, err)
		runCheck = true
	}
	fakeEndpoint.SendRequestHandleResponse(ctx, node0, req, resp, respHandler)
	require.Equal(t, endpoint.ErrUnknownPeer, err)

	err = fakeEndpoint.MaybeAddToPeerstore(ctx, node0, peerstoreTTL)
	require.NoError(t, err)

	connectedness = fakeEndpoint.Connectedness(node0)
	require.Equal(t, network.CanConnect, connectedness)

	na, err := fakeEndpoint.NetworkAddress(node0)
	require.NoError(t, err)
	require.Equal(t, node0, na)

	fakeEndpoint.SendRequestHandleResponse(ctx, node0, req, resp, respHandler)

	sched0 := simplescheduler.NewSimpleScheduler(ctx, clk)
	rt0 := simplert.NewSimpleRT(kadid.Key(), 2)
	serv0 := kadsimserver.NewKadSimServer(rt0, fakeEndpoint)
	dispatcher.AddPeer(node0, sched0, serv0)

	fakeEndpoint.SendRequestHandleResponse(ctx, node0, req, resp, respHandler)

	require.True(t, sched0.RunOne(ctx))
	require.False(t, sched0.RunOne(ctx))

	require.True(t, runCheck)
}
