package query

import (
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func TestAddPeers(t *testing.T) {
	// create empty peer list
	pl := newPeerList(key.ZeroKey)

	require.Nil(t, pl.closest)
	require.Nil(t, pl.closestQueued)

	// add initial peers
	nPeers := 3
	peerids := make([]peer.ID, nPeers+1)
	for i := 0; i < nPeers; i++ {
		peerids[i] = peer.ID(byte(i))
	}
	peerids[nPeers] = peer.ID(byte(0)) // duplicate with peerids[0]

	// distances
	// peerids[0]: 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d
	// peerids[1]: 4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a
	// peerids[2]: dbc1b4c900ffe48d575b5da5c638040125f65db0fe3e24494b76ea986457d986
	// peerids[3]: 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d

	// add 4 peers (incl. 1 duplicate)
	addToPeerlist(pl, peerids)

	curr := pl.closest
	// verify that closest peer is peerids[1]
	require.Equal(t, peerids[1], curr.id)
	curr = curr.next
	// second closest peer should be peerids[0]
	require.Equal(t, peerids[0], curr.id)
	curr = curr.next
	// third closest peer should be peerids[2]
	require.Equal(t, peerids[2], curr.id)

	// end of the list
	require.Nil(t, curr.next)

	// verify that closestQueued peer is peerids[0]
	require.Equal(t, peerids[1], pl.closestQueued.id)

	// add more peers
	nPeers = 5
	newPeerids := make([]peer.ID, nPeers+2)
	for i := 0; i < nPeers; i++ {
		newPeerids[i] = peer.ID(byte(10 + i))
	}
	newPeerids[nPeers] = peer.ID(byte(10))  // duplicate with newPeerids[0]
	newPeerids[nPeers+1] = peer.ID(byte(1)) // duplicate with peerids[1]

	// distances
	// newPeerids[0]: 01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b
	// newPeerids[1]: e7cf46a078fed4fafd0b5e3aff144802b853f8ae459a4f0c14add3314b7cc3a6
	// newPeerids[2]: ef6cbd2161eaea7943ce8693b9824d23d1793ffb1c0fca05b600d3899b44c977
	// newPeerids[3]: 9d1e0e2d9459d06523ad13e28a4093c2316baafe7aec5b25f30eba2e113599c4
	// newPeerids[4]: 4d7b3ef7300acf70c892d8327db8272f54434adbc61a4e130a563cb59a0d0f47
	// newPeerids[5]: 01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b
	// newPeerids[6]: 4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a

	// add 7 peers (incl. 2 duplicates)
	addToPeerlist(pl, newPeerids)

	// order is now as follows:
	order := []peer.ID{newPeerids[0], peerids[1], newPeerids[4], peerids[0], newPeerids[3],
		peerids[2], newPeerids[1], newPeerids[2]}

	curr = pl.closest
	for _, p := range order {
		require.Equal(t, p, curr.id)
		curr = curr.next
	}
	require.Nil(t, curr)

	// verify that closestQueued peer is peerids[0]
	require.Equal(t, newPeerids[0], pl.closestQueued.id)

	// add a single peer that isn't the closest one
	newPeer := peer.ID(byte(20))

	addToPeerlist(pl, []peer.ID{newPeer})
	order = append(order[:5], order[4:]...)
	order[4] = newPeer

	curr = pl.closest
	for _, p := range order {
		require.Equal(t, p, curr.id)
		curr = curr.next
	}

	require.Nil(t, curr)
}

func TestPeerlistCornerCases(t *testing.T) {
	// different empty peerid lists
	emptyPeeridLists := [][]peer.ID{
		{},
		{peer.ID("")},
		make([]peer.ID, 3),
	}

	for _, peerids := range emptyPeeridLists {
		// adding them to a peerlist should result in an empty list
		require.Nil(t, sliceToPeerInfos(key.ZeroKey, peerids))

		pl := newPeerList(key.ZeroKey)
		addToPeerlist(pl, peerids)
		require.Nil(t, pl.closest)
		require.Nil(t, pl.closestQueued)
	}

	pl := newPeerList(key.ZeroKey)

	singlePeerList0 := []peer.ID{peer.ID(byte(0))}
	// 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d
	addToPeerlist(pl, singlePeerList0)
	require.Equal(t, singlePeerList0[0], pl.closest.id)
	require.Equal(t, singlePeerList0[0], pl.closestQueued.id)

	singlePeerList1 := []peer.ID{peer.ID(byte(1))}
	// 4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a
	addToPeerlist(pl, singlePeerList1)
	require.Equal(t, singlePeerList1[0], pl.closest.id)
	require.Equal(t, singlePeerList1[0], pl.closest.id)

	singlePeerList2 := []peer.ID{peer.ID(byte(2))}
	// dbc1b4c900ffe48d575b5da5c638040125f65db0fe3e24494b76ea986457d986
	addToPeerlist(pl, singlePeerList2)
	curr := pl.closest
	require.Equal(t, singlePeerList1[0], curr.id)
	curr = curr.next
	require.Equal(t, singlePeerList0[0], curr.id)
	curr = curr.next
	require.Equal(t, singlePeerList2[0], curr.id)
	curr = curr.next
	require.Nil(t, curr)
}

func TestUpdatePeerStatusInPeerlist(t *testing.T) {
	// create empty peer list
	pl := newPeerList(key.ZeroKey)

	// add initial peers
	nPeers := 3
	peerids := make([]peer.ID, nPeers)
	for i := 0; i < nPeers; i++ {
		peerids[i] = peer.ID(byte(i))
	}

	addToPeerlist(pl, peerids)

	// initial queue state
	// peerids[1], queued, 4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a
	// peerids[0], queued, 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d
	// peerids[2], queued, dbc1b4c900ffe48d575b5da5c638040125f65db0fe3e24494b76ea986457d986

	updatePeerStatusInPeerlist(pl, peerids[0], waiting)
	updatePeerStatusInPeerlist(pl, peerids[1], unreachable)

	curr := pl.closest
	require.Equal(t, curr.status, unreachable)
	curr = curr.next
	require.Equal(t, curr.status, waiting)
	curr = curr.next
	require.Equal(t, curr.status, queued)
	curr = curr.next
	require.Nil(t, curr)
}

func TestPopClosestQueued(t *testing.T) {
	// create empty peer list
	pl := newPeerList(key.ZeroKey)

	// add initial peers
	nPeers := 3
	peerids := make([]peer.ID, nPeers)
	for i := 0; i < nPeers; i++ {
		peerids[i] = peer.ID(byte(i))
	}

	addToPeerlist(pl, peerids)

	// initial queue state
	// peerids[1], queued, 4bf5122f344554c53bde2ebb8cd2b7e3d1600ad631c385a5d7cce23c7785459a
	// peerids[0], queued, 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d
	// peerids[2], queued, dbc1b4c900ffe48d575b5da5c638040125f65db0fe3e24494b76ea986457d986

	require.Equal(t, peerids[1], popClosestQueued(pl))
	require.Equal(t, peerids[1], pl.closest.id)
	require.Equal(t, peerids[0], pl.closestQueued.id)
	require.Equal(t, peerids[0], popClosestQueued(pl))
	require.Equal(t, peerids[1], pl.closest.id)
	require.Equal(t, peerids[2], pl.closestQueued.id)
	require.Equal(t, peerids[2], popClosestQueued(pl))
	require.Equal(t, peer.ID(""), popClosestQueued(pl))

	pl = newPeerList(key.ZeroKey)

	addToPeerlist(pl, peerids)

	// mark second item (peerids[0]) as waiting
	updatePeerStatusInPeerlist(pl, peerids[0], waiting)
	require.Equal(t, peerids[1], pl.closest.id)
	require.Equal(t, peerids[1], pl.closestQueued.id)

	// pop closest queued (peerids[1])
	require.Equal(t, peerids[1], popClosestQueued(pl))
	require.Equal(t, peerids[1], pl.closest.id)
	// peerids[2] is now closestQueued
	require.Equal(t, peerids[2], pl.closestQueued.id)

	updatePeerStatusInPeerlist(pl, peerids[2], unreachable)
	require.Equal(t, peerids[1], pl.closest.id)
	require.Nil(t, pl.closestQueued)
}
