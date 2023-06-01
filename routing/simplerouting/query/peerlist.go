package query

import (
	"sort"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p/core/peer"
)

type peerStatus uint8

const (
	queued peerStatus = iota
	waiting
	queried
	unreachable
)

type peerInfo struct {
	distance key.KadKey
	status   peerStatus
	id       peer.ID

	next *peerInfo
}

type peerList struct {
	target key.KadKey

	closest       *peerInfo
	closestQueued *peerInfo
}

func newPeerList(target key.KadKey) *peerList {
	return &peerList{
		target: target,
	}
}

// normally peers should already be ordered with distance to target, but we
// sort them just in case
// TODO: check corner cases with len(peers) == 0 or 1
func addToPeerlist(pl *peerList, peers []peer.ID) {

	// linked list of new peers sorted by distance to target
	newHead := sliceToPeerInfos(pl.target, peers)

	// if the list is empty, define first new peer as closest
	if pl.closest == nil {
		pl.closest = newHead
		pl.closestQueued = newHead
		return
	}

	// merge the new sorted list into the existing sorted list
	var prev *peerInfo
	currOld := true
	closestQueuedReached := false

	oldHead := pl.closest

	// TODO: update closestQueued pointer

	r := key.Compare(oldHead.distance, newHead.distance)
	if r > 0 {
		pl.closest = newHead
		pl.closestQueued = newHead
		currOld = false
	}

	for {
		if r > 0 {
			// newHead is closer than oldHead

			if !closestQueuedReached {
				// newHead is closer than closestQueued, update closestQueued
				pl.closestQueued = newHead
				closestQueuedReached = true
			}
			if currOld && prev != nil {
				prev.next = newHead
				currOld = false
			}
			prev = newHead
			newHead = newHead.next
		} else {
			// oldHead is closer than newHead

			if !closestQueuedReached && oldHead == pl.closestQueued {
				// old closestQueued is closer than newHead,
				// don't update closestQueued
				closestQueuedReached = true
			}

			if !currOld && prev != nil {
				prev.next = oldHead
				currOld = true
			}
			prev = oldHead
			oldHead = oldHead.next
			if r == 0 {
				// newHead is a duplicate of oldHead, discard newHead
				newHead = newHead.next
			}
		}
		// we are done when we reach the end of either list
		if oldHead == nil || newHead == nil {
			break
		}
		r = key.Compare(oldHead.distance, newHead.distance)
	}

	// append the remaining list to the end
	if oldHead == nil {
		prev.next = newHead
	} else {
		prev.next = oldHead
	}

}

func sliceToPeerInfos(target key.KadKey, peers []peer.ID) *peerInfo {
	if len(peers) == 0 {
		return nil
	}

	// create a new list of peerInfo
	newPeers := make([]peerInfo, len(peers))
	for i, p := range peers {
		newPeers[i] = addrInfoToPeerInfo(target, p)
	}

	// sort the new list
	sort.Slice(newPeers, func(i, j int) bool {
		return key.Compare(newPeers[i].distance, newPeers[j].distance) < 0
	})

	// convert slice to linked list and remove duplicates
	curr := &newPeers[0]
	for i := 1; i < len(newPeers); i++ {
		if curr.distance != newPeers[i].distance {
			curr.next = &newPeers[i]
			curr = curr.next
		}
	}
	// return head of linked list
	return &newPeers[0]
}

func addrInfoToPeerInfo(target key.KadKey, p peer.ID) peerInfo {
	return peerInfo{
		distance: key.Xor(target, key.PeerKadID(p)),
		status:   queued,
		id:       p,
	}
}

func popClosestQueued(pl *peerList) peer.ID {
	if pl.closestQueued == nil {
		return peer.ID("")
	}
	pi := pl.closestQueued
	pi.status = waiting
	curr := pl.closestQueued
	for curr.next != nil && curr.next.status != queued {
		curr = curr.next
	}
	pl.closestQueued = curr
	return pi.id
}

func updatePeerStatusInPeerlist(pl *peerList, p peer.ID, newStatus peerStatus) {
	curr := pl.closest
	for curr != nil && curr.id != p {
		curr = curr.next
	}
	if curr != nil {
		curr.status = newStatus
	}
}
