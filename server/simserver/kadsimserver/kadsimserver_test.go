package kadsimserver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/network/address/kadid"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint/fakeendpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/simmessage"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

var self = kadid.KadID{KadKey: []byte{0x00}} // 0000 0000

// remotePeers with bucket assignments wrt to self
var remotePeers = []*kadid.KadID{
	{KadKey: []byte{0b10001000}}, // 1000 1000 (bucket 0)
	{KadKey: []byte{0b11010010}}, // 1101 0010 (bucket 0)
	{KadKey: []byte{0b01001011}}, // 0100 1011 (bucket 1)
	{KadKey: []byte{0b01010011}}, // 0101 0011 (bucket 1)
	{KadKey: []byte{0b00101110}}, // 0010 1110 (bucket 2)
	{KadKey: []byte{0b00110110}}, // 0011 0110 (bucket 2)
	{KadKey: []byte{0b00011111}}, // 0001 1111 (bucket 3)
	{KadKey: []byte{0b00010001}}, // 0001 0001 (bucket 3)
	{KadKey: []byte{0b00001000}}, // 0000 1000 (bucket 4)
}

func TestKadSimServer(t *testing.T) {
	ctx := context.Background()

	peerstoreTTL := time.Second // doesn't matter as we use fakeendpoint
	numberOfCloserPeersToSend := 4

	fakeEndpoint := fakeendpoint.NewFakeEndpoint(self, nil)
	rt := simplert.NewSimpleRT(self.Key(), 2)

	// add peers to routing table and peerstore
	for _, p := range remotePeers {
		err := fakeEndpoint.MaybeAddToPeerstore(ctx, p, peerstoreTTL)
		require.NoError(t, err)
		success, err := rt.AddPeer(ctx, p)
		require.NoError(t, err)
		require.True(t, success)
	}

	s0 := NewKadSimServer(rt, fakeEndpoint, WithPeerstoreTTL(peerstoreTTL),
		WithNumberOfCloserPeersToSend(numberOfCloserPeersToSend))
	var runCount int

	requester := kadid.KadID{KadKey: []byte{0b00000001}} // 0000 0001

	req0 := simmessage.NewSimRequest([]byte{0b00000000})
	check0 := func(resp message.MinKadResponseMessage) {
		require.Len(t, resp.CloserNodes(), numberOfCloserPeersToSend)
		// closer peers should be ordered by distance to 0000 0000
		// [8] 0000 1000, [7] 0001 0001, [6] 0001 1111, [4] 0010 1110
		order := []*kadid.KadID{remotePeers[8], remotePeers[7], remotePeers[6], remotePeers[4]}
		for i, p := range resp.CloserNodes() {
			require.Equal(t, order[i], p)
		}
		runCount++
	}
	s0.HandleFindNodeRequest(ctx, requester, req0, check0)
	require.Equal(t, 1, runCount)

	req1 := simmessage.NewSimRequest([]byte{0b11111111})
	check1 := func(resp message.MinKadResponseMessage) {
		require.Len(t, resp.CloserNodes(), numberOfCloserPeersToSend)
		// closer peers should be ordered by distance to 1111 1111
		// [1] 1101 0010, [0] 1000 1000, [3] 0101 0011, [2] 0100 1011
		order := []*kadid.KadID{remotePeers[1], remotePeers[0], remotePeers[3], remotePeers[2]}
		for i, p := range resp.CloserNodes() {
			require.Equal(t, order[i], p)
		}
		runCount++
	}
	s0.HandleFindNodeRequest(ctx, requester, req1, check1)
	require.Equal(t, 2, runCount)

	numberOfCloserPeersToSend = 3
	s1 := NewKadSimServer(rt, fakeEndpoint, WithNumberOfCloserPeersToSend(3))
	runCount = 0

	req2 := simmessage.NewSimRequest([]byte{0b01100000})
	check2 := func(resp message.MinKadResponseMessage) {
		require.Len(t, resp.CloserNodes(), numberOfCloserPeersToSend)
		// closer peers should be ordered by distance to 0110 0000
		// [2] 0100 1011, [3] 0101 0011, [4] 0010 1110
		order := []*kadid.KadID{remotePeers[2], remotePeers[3], remotePeers[4]}
		for i, p := range resp.CloserNodes() {
			require.Equal(t, order[i], p)
		}
		runCount++
	}
	s1.HandleFindNodeRequest(ctx, requester, req2, check2)
	require.Equal(t, 1, runCount)
}

func TestInvalidRequests(t *testing.T) {
	ctx := context.Background()
	// invalid option
	s := NewKadSimServer(nil, nil, func(*Config) error {
		return errors.New("invalid option")
	})
	require.Nil(t, s)

	peerstoreTTL := time.Second // doesn't matter as we use fakeendpoint

	// create a valid server
	fakeEndpoint := fakeendpoint.NewFakeEndpoint(self, nil)
	rt := simplert.NewSimpleRT(self.Key(), 2)

	// add peers to routing table and peerstore
	for _, p := range remotePeers {
		err := fakeEndpoint.MaybeAddToPeerstore(ctx, p, peerstoreTTL)
		require.NoError(t, err)
		success, err := rt.AddPeer(ctx, p)
		require.NoError(t, err)
		require.True(t, success)
	}

	s = NewKadSimServer(rt, fakeEndpoint)
	require.NotNil(t, s)

	requester := kadid.KadID{KadKey: []byte{0b00000001}} // 0000 0001

	p, err := peer.Decode("1D3oooVaLidPeerid")
	require.NoError(t, err)
	pid := peerid.NewPeerID(p)

	// invalid message format (not a SimMessage)
	req0 := ipfskadv1.FindPeerRequest(pid)
	check := func(resp message.MinKadResponseMessage) {
		require.Fail(t, "response function should not be called")
	}
	s.HandleFindNodeRequest(ctx, requester, req0, check)

	// empty request
	req1 := &simmessage.SimMessage{}
	s.HandleFindNodeRequest(ctx, requester, req1, check)

	// request with invalid key (not matching the expected length)
	req2 := simmessage.NewSimRequest([]byte{0b00000000, 0b00000001})
	s.HandleFindNodeRequest(ctx, requester, req2, check)
}
