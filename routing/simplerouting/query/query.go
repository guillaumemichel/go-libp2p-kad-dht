package query

import (
	"context"
	"errors"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	// MAGIC: takes the default value from the DHT constants
	NClosestPeers = consts.NClosestPeers
)

var (
	// MAGIC: default peerstore TTL for newly discovered peers
	QueuedPeersPeerstoreTTL = peerstore.TempAddrTTL
)

type SimpleQuery struct {
	ctx         context.Context
	done        bool
	kadid       key.KadKey
	message     *pb.Message
	concurrency int
	timeout     time.Duration
	proto       protocol.ID

	host        host.Host // TODO: for now only peerstore required
	msgEndpoint *network.MessageEndpoint
	rt          routingtable.RoutingTable

	eventqueue eq.EventQueue
	sched      events.Scheduler

	peerlist *peerList

	// success condition
	interestingElements []interface{}
	resultsChan         chan interface{}
	testSuccess         func(context.Context, []interface{}, *pb.Message, chan interface{}) (bool, []interface{})
}

func NewSimpleQuery(ctx context.Context, kadid key.KadKey, message *pb.Message,
	concurrency int, timeout time.Duration, proto protocol.ID, h host.Host,
	msgEndpoint *network.MessageEndpoint, rt routingtable.RoutingTable,
	queue eq.EventQueue, sched events.Scheduler, resultsChan chan interface{},
	success func(context.Context, []interface{}, *pb.Message,
		chan interface{}) (bool, []interface{})) *SimpleQuery {

	ctx, span := internal.StartSpan(ctx, "SimpleQuery.NewSimpleQuery",
		trace.WithAttributes(attribute.String("Target", kadid.Hex())))
	defer span.End()

	closestPeers := rt.NearestPeers(ctx, kadid, NClosestPeers)

	peerlist := newPeerList(kadid)
	addToPeerlist(peerlist, closestPeers)

	query := &SimpleQuery{
		ctx:                 ctx,
		message:             message,
		kadid:               kadid,
		concurrency:         concurrency,
		timeout:             timeout,
		proto:               proto,
		host:                h,
		msgEndpoint:         msgEndpoint,
		rt:                  rt,
		peerlist:            peerlist,
		eventqueue:          queue,
		sched:               sched,
		interestingElements: []interface{}{},
		resultsChan:         resultsChan,
		testSuccess:         success,
	}

	// TODO: add concurrency request events to eventqueue
	for i := 0; i < concurrency; i++ {
		query.eventqueue.Enqueue(query.newRequest)
	}

	return query
}

func (q *SimpleQuery) checkIfDone() error {
	if q.done {
		// query is done, don't send any more requests
		return errors.New("query done")
	}

	select {
	case <-q.ctx.Done():
		// query is cancelled, mark it as done
		q.done = true
		return errors.New("query cancelled")
	default:
	}
	return nil
}

func (q *SimpleQuery) newRequest() {
	ctx, span := internal.StartSpan(q.ctx, "SimpleQuery.newRequest")
	defer span.End()

	if err := q.checkIfDone(); err != nil {
		span.AddEvent(err.Error())
		return
	}

	peerid := popClosestQueued(q.peerlist)
	if peerid == "" {
		// TODO: handle this case
		span.AddEvent("all peers queried")
		return
	}
	span.AddEvent("peer selected: " + peerid.String())

	// start new go routine to send request to peer
	go q.sendRequest(ctx, peerid)

	// add timeout to scheduler
	events.ScheduleAction(ctx, &q.sched, q.timeout, func() {
		q.requestError(peerid, errors.New("request timeout ("+q.timeout.String()+")"))
	})
}

func (q *SimpleQuery) sendRequest(ctx context.Context, p peer.ID) {
	ctx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()

	ctx, span := internal.StartSpan(ctx, "SimpleQuery.sendRequest")
	defer span.End()

	if err := q.msgEndpoint.DialPeer(ctx, p); err != nil {
		span.AddEvent("peer dial failed")
		q.eventqueue.Enqueue(func() {
			q.requestError(p, err)
		})
		return
	}

	resp, err := q.msgEndpoint.SendRequest(ctx, p, q.message, q.proto)
	if err != nil {
		span.AddEvent("request failed")
		q.eventqueue.Enqueue(func() {
			q.requestError(p, err)
		})
		return
	}
	span.AddEvent("got a response")
	q.eventqueue.Enqueue(func() {
		q.handleResponse(p, resp)
	})
}

func (q *SimpleQuery) handleResponse(p peer.ID, resp *pb.Message) {
	ctx, span := internal.StartSpan(q.ctx, "SimpleQuery.handleResponse")
	defer span.End()

	if err := q.checkIfDone(); err != nil {
		span.AddEvent(err.Error())
		return
	}

	closerPeers := resp.GetCloserPeers()
	if len(closerPeers) > 0 {
		// consider that remote peer is behaving correctly if it returns
		// at least 1 peer
		q.rt.AddPeer(ctx, p)
	}

	var success bool
	success, q.interestingElements = q.testSuccess(ctx, q.interestingElements, resp, q.resultsChan)
	if success {
		// query is done, don't send any more requests
		span.AddEvent("query success")
		q.done = true
		return
	}

	updatePeerStatusInPeerlist(q.peerlist, p, queried)

	newPeers := network.ParsePeers(ctx, closerPeers)
	newPeerIds := make([]peer.ID, 0, len(newPeers))

	for _, ai := range newPeers {
		if ai.ID == q.msgEndpoint.Host.ID() {
			// don't add self to queries or routing table
			span.AddEvent("remote peer provided self as closer peer")
			continue
		}
		newPeerIds = append(newPeerIds, ai.ID)

		q.msgEndpoint.MaybeAddToPeerstore(ai, QueuedPeersPeerstoreTTL)
	}

	addToPeerlist(q.peerlist, newPeerIds)

	// add pending request for this query to eventqueue
	q.eventqueue.Enqueue(q.newRequest)
}

func (q *SimpleQuery) requestError(peerid peer.ID, err error) {
	ctx, span := internal.StartSpan(q.ctx, "SimpleQuery.requestError",
		trace.WithAttributes(attribute.String("PeerID", peerid.String()),
			attribute.String("Error", err.Error())))
	defer span.End()

	if q.ctx.Err() == nil {
		// remove peer from routing table unless context was cancelled
		q.rt.RemovePeer(ctx, key.PeerKadID(peerid))
	}

	if err := q.checkIfDone(); err != nil {
		span.AddEvent(err.Error())
		return
	}

	updatePeerStatusInPeerlist(q.peerlist, peerid, unreachable)

	// add pending request for this query to eventqueue
	q.eventqueue.Enqueue(q.newRequest)
}

func (q *SimpleQuery) Close() {
	// TODO: check if we need to cancel anything else
	q.done = true
	q.resultsChan <- errors.New("query cancelled")
	close(q.resultsChan)
}
