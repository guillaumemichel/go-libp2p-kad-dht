package routing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p/core/peer"
)

const PEERSTORE_ENTRY_LIFETIME = 30 * time.Minute

type query struct {
	// query peer set
	// {dist, conn_type, peer.Info, queried?}

	ctx   context.Context
	kadId hash.KadKey

	req                  *pb.DhtMessage
	lk                   sync.Mutex
	done                 bool
	interestingResponses []*pb.DhtMessage

	qpeerset *qpeerset
	results  chan queryResult

	routing *DhtRouting
}

type queryResult struct {
	// errors
	peer peer.AddrInfo
	// values
}

type queryManager struct {
	limit          int
	ongoingQueries []*query
	newQueries     chan *query

	kill chan struct{}

	routing *DhtRouting
}

func (r *DhtRouting) newQueryManager(limit int) *queryManager {
	qm := &queryManager{
		limit:          limit,
		ongoingQueries: []*query{},
		newQueries:     make(chan *query),
		routing:        r,
	}
	qm.run()
	return qm
}

func (qm *queryManager) run() {
	for {
		if len(qm.ongoingQueries) < qm.limit {
			if len(qm.queuedQueries) > 0 {
				q := qm.queuedQueries[0]
				qm.queuedQueries = qm.queuedQueries[1:]
				qm.ongoingQueries = append(qm.ongoingQueries, q)
				go q.run()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// can return either a peer or a value (or both).
// peer is peerid + multiaddrs or peerid only (peer.AddrInfo)
// value is a byte array (or string)

func (qm *queryManager) newQuery(ctx context.Context, kadId hash.KadKey) *query {
	return &query{
		ctx:                  ctx,
		kadId:                kadId,
		qpeerset:             newQpeerset(kadId),
		routing:              qm.routing,
		results:              make(chan queryResult),
		interestingResponses: []*pb.DhtMessage{},
	}
}

func (qm *queryManager) Query(ctx context.Context, kadId hash.KadKey, msg *pb.DhtMessage, stopCond func([]*pb.DhtMessage) (bool, []*pb.DhtMessage)) chan queryResult {
	q := qm.newQuery(ctx, kadId)
	qm.queuedQueries = append(qm.queuedQueries, q)
	return q.results
}

func (q *query) run() {
	select {
	case <-q.ctx.Done():
		q.results <- queryResult{} // error
		return
	default:
	}

	localClosest := q.routing.rt.NearestPeers(q.kadId, q.routing.concurrency)
	for _, p := range localClosest {
		q.qpeerset.AddPeer(p)
	}

	for i := 0; i < q.routing.concurrency; i++ {
		go func() {
			for {
				// until ctx is cancelled, or we have enough results, or we don't learn about new peers
				q.lk.Lock()
				qpeers := q.qpeerset.ClosestHeard(1)
				q.lk.Unlock()
				resp, err := q.routing.me.SendDhtRequest(q.ctx, qpeers[0].ID, q.req)
				if err != nil {
					fmt.Println("error sending request: ", err)
					continue
				}
				err = q.handleResponse(resp)
				if err != nil {
					fmt.Println("error handling response: ", err)
					continue
				}
			}
		}()
	}
}

func (q *query) handleResponse(resp *pb.DhtMessage) error {
	if resp.GetFindPeerResponseType() != nil {
		return q.handleFindPeerResponse(resp.GetFindPeerResponseType())
	} else {
		return fmt.Errorf("unexpected response type: %v", resp)
	}
}

func (q *query) handleFindPeerResponse(resp *pb.DhtFindPeerResponse) error {
	q.lk.Lock()
	defer q.lk.Unlock()
	if q.done { // do something with ctx
		return fmt.Errorf("query cancelled")
	}

	peers, parseErr := dhtnet.PBPeerToPeerInfos(resp.GetPeers())
	if parseErr != nil {
		fmt.Println("error parsing peers: ", parseErr)
	}

	qpeers := q.handlePeers(peers)
	for _, qpeer := range qpeers {
		if qpeer.dist.Compare(hash.ZeroKey) == 0 {
			// found the peer we were looking for
			q.done = true
			q.results <- queryResult{peer: qpeer.AddrInfo}
		}
	}

	if len(peers) > 0 && len(qpeers) == 0 {
		// no new peers
		q.done = true
		return fmt.Errorf("no new peers")
	}
	return nil
}

func (q *query) handlePeers(peers []peer.AddrInfo) []*qpeer {

	newPeers := make([]*qpeer, 0, len(peers))
	// add to query peer set
	q.lk.Lock()
	for _, p := range peers {
		newPeers = append(newPeers, q.qpeerset.AddPeer(p))
	}
	q.lk.Unlock()

	for _, p := range peers {
		// add to peerstore
		q.routing.me.Host.Peerstore().AddAddrs(p.ID, p.Addrs, PEERSTORE_ENTRY_LIFETIME)
		// add to routing table
		q.routing.rt.AddPeer(p)
	}

	return newPeers
}

type peerState uint8

const (
	// PeerHeard is applied to peers which have not been queried yet.
	PeerHeard peerState = iota
	// PeerWaiting is applied to peers that are currently being queried.
	PeerWaiting
	// PeerQueried is applied to peers who have been queried and a response was retrieved successfully.
	PeerQueried
	// PeerUnreachable is applied to peers who have been queried and a response was not retrieved successfully.
	PeerUnreachable
)

type qpeer struct {
	peer.AddrInfo
	dist      hash.KadKey
	peerState peerState

	next *qpeer
}

type qpeerset struct {
	key hash.KadKey

	head *qpeer
	size int
}

func newQpeerset(key hash.KadKey) *qpeerset {
	return &qpeerset{
		key: key,
	}
}

func (qps *qpeerset) AddPeer(ai peer.AddrInfo) *qpeer {
	qpeer := &qpeer{
		AddrInfo:  ai,
		dist:      qps.key.Xor(hash.PeerKadID(ai.ID)),
		peerState: PeerHeard,
	}
	if qps.insert(qpeer) {
		return qpeer
	}
	return nil
}

// insert inserts a peer from the top
// doesn't insert duplicates
func (qps *qpeerset) insert(p *qpeer) bool {
	if qps.head == nil {
		qps.head = p
		qps.size++
		return true
	}

	var prev *qpeer
	curr := qps.head

	for curr != nil && p.dist.Compare(curr.dist) > 0 {
		prev = curr
		curr = curr.next
	}

	if curr != nil && p.dist.Compare(curr.dist) == 0 {
		// duplicate
		return false
	}

	if prev == nil {
		// insert at head
		p.next = qps.head
		qps.head = p
	} else {
		// insert between prev and curr
		p.next = prev.next
		prev.next = p
	}

	qps.size++
	return true
}

// Closest Heard returns the n closest peers that have been heard of, but not
// queried yet.
func (qps *qpeerset) ClosestHeard(n int) []*qpeer {
	heard := make([]*qpeer, 0, n)
	curr := qps.head
	for curr != nil && n > 0 {
		if curr.peerState == PeerHeard {
			heard = append(heard, curr)
			n--
		}
		curr = curr.next
	}
	return heard
}
