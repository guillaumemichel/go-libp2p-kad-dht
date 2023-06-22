package ipfssimserver

/*
func TestSimIpfsServer(t *testing.T) {
	ctx := context.Background()

	peerstoreTTL := time.Second // doesn't matter as we use fakeendpoint
	numberOfCloserPeersToSend := 4

	selfPid, err := peer.Decode("1ECoooSELF")
	// KadKey: b761a0e5365eabea35f325f1b9f8d27c32d4395d35a6a3dad93ceab4f5553628
	require.NoError(t, err)
	self := peerid.NewPeerID(selfPid)

	nRemotePeers := 7
	remotePeers := make([]*peerid.PeerID, nRemotePeers)
	// generate peerids for remote peers (1EoooPEER2, 1EoooPEER3, ..., 1EoooPEER8)
	for i := 0; i < nRemotePeers; i++ {
		pid, err := peer.Decode(fmt.Sprintf("1EoooPEER%d", i+2))
		require.NoError(t, err, i)
		remotePeers[i] = peerid.NewPeerID(pid)
		fmt.Println("	// ["+strconv.Itoa(i)+"]", remotePeers[i].Key())
	}
	// remote peers
	// [0] e69614c5fcb92e8fbf2aa5785904fec5a67524ac7cd513f32bc7ab38621b4b7b (bucket 1)
	// [1] 6ab9cb73bbd52ad2bb6ac4048e988478bf076df9b39e072f30b4722639382683 (bucket 0)
	// [2] 69b9104f74ca05073a1bb658155fa4549fcc8db470947915a6e2750185dc1f81 (bucket 0)
	// [3] 4eaafc67b177fa53ee6de27d1646f7862fb2957878bcc8d60dfa67b7832bb28b (bucket 0)
	// [4] ab6c9fe862d32ff3170ed43600742b2abbb52f09216afa139cb89842e083ce4e (bucket 3)
	// [5] 00ca8d64555add66790c4fb3e62075911a02a3577622fa69279731e82c135b8a (bucket 0)
	// [6] 7bc87a73b9223cd1c936f393f3abd9e2246946195382bb719c49cb44a4e9afb4 (bucket 0)

	fakeEndpoint := fakeendpoint.NewFakeEndpoint(self, nil)
	rt := simplert.NewSimpleRT(self.Key(), 5)

	// add peers to routing table and peerstore
	for _, p := range remotePeers {
		err := fakeEndpoint.MaybeAddToPeerstore(ctx, p, peerstoreTTL)
		require.NoError(t, err)
		success, err := rt.AddPeer(ctx, p)
		require.NoError(t, err)
		require.True(t, success)
	}

	s0 := NewIpfsSimServer(rt, fakeEndpoint, WithPeerstoreTTL(peerstoreTTL),
		WithNumberOfCloserPeersToSend(numberOfCloserPeersToSend))
	var runCount int

	reqPid, err := peer.Decode("1WoooREQUESTER")
	require.NoError(t, err)
	requester := peerid.NewPeerID(reqPid)

	targetPid, err := peer.Decode("12BoooTARGET")
	require.NoError(t, err)
	target := peerid.NewPeerID(targetPid)

	req0 := ipfskadv1.FindPeerRequest(target)
	check0 := func(resp message.MinKadResponseMessage) {
		require.Equal(t, requester, resp.GetRequester())
		require.Equal(t, self, resp.GetResponder())
		require.Equal(t, []byte{0b00000000}, resp.GetKey())
		require.Equal(t, []byte{0b00000000}, resp.GetValue())
		require.Equal(t, []byte{0b00000000}, resp.GetCloserPeers())
		require.Equal(t, []byte{0b00000000}, resp.GetProviderPeers())
		require.Equal(t, []byte{0b00000000}, resp.GetExtra())
	}

	// run 0
	runCount++
	resp0, err := s0.HandleSimRequest(ctx, req0)
	require.NoError(t, err)
	check0(resp0)

	// run 1
	runCount++
	s1 := NewKadSimServer(rt, fakeEndpoint, WithPeerstoreTTL(peerstoreTTL),
		WithNumberOfCloserPeersToSend(numberOfCloserPeersToSend))
	resp1, err := s1.HandleSimRequest(ctx, req0)
	require.NoError(t, err)
	check0(resp1)

	// run 2
	runCount++
	s2 := NewKadSimServer(rt, fakeEndpoint, WithPeerstoreTTL(peerstoreTTL),
		WithNumberOfCloserPeersToSend(numberOfCloserPeersToSend))
	resp2, err := s2.HandleSimRequest(ctx, req0)
	require.NoError(t, err)
}
*/
