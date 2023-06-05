package simplequery

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dht/consts"
	"github.com/libp2p/go-libp2p-kad-dht/events"
	eq "github.com/libp2p/go-libp2p-kad-dht/events/eventqueue"
	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
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

type HandleResultFn func(context.Context, []interface{}, *pb.Message, chan interface{}) []interface{}

type SimpleQuery struct {
	ctx         context.Context
	self        key.KadKey
	done        bool
	kadid       key.KadKey
	message     *pb.Message
	concurrency int
	timeout     time.Duration
	proto       protocol.ID

	msgEndpoint endpoint.Endpoint
	rt          routingtable.RoutingTable

	eventQueue   eq.EventQueue
	eventPlanner events.EventPlanner

	inflightRequests int // requests that are either in flight or scheduled
	peerlist         *peerList

	// success condition
	state          []interface{}
	resultsChan    chan interface{}
	handleResultFn HandleResultFn
}

// NewSimpleQuery creates a new SimpleQuery. It initializes the query by adding
// the closest peers to the target key from the provided routing table to the
// query's peerlist. It sends `concurreny` requests events to the provided event
// queue. The requests events and followup events are handled by the event queue
// reader, and the parameters to these events are determined by the query's
// parameters. The query keeps track of the closest known peers to the target
// key, and the peers that have been queried so far.
func NewSimpleQuery(ctx context.Context, self, kadid key.KadKey, message *pb.Message,
	concurrency int, timeout time.Duration, proto protocol.ID,
	msgEndpoint endpoint.Endpoint, rt routingtable.RoutingTable,
	queue eq.EventQueue, ep events.EventPlanner, resultsChan chan interface{},
	handleResultFn HandleResultFn) *SimpleQuery {

	ctx, span := internal.StartSpan(ctx, "SimpleQuery.NewSimpleQuery",
		trace.WithAttributes(attribute.String("Target", kadid.Hex())))
	defer span.End()

	closestPeers := rt.NearestPeers(ctx, kadid, NClosestPeers)

	peerlist := newPeerList(kadid)
	addToPeerlist(peerlist, closestPeers)

	query := &SimpleQuery{
		ctx:              ctx,
		self:             self,
		message:          message,
		kadid:            kadid,
		concurrency:      concurrency,
		timeout:          timeout,
		proto:            proto,
		msgEndpoint:      msgEndpoint,
		rt:               rt,
		inflightRequests: 0,
		peerlist:         peerlist,
		eventQueue:       queue,
		eventPlanner:     ep,
		state:            []interface{}{},
		resultsChan:      resultsChan,
		handleResultFn:   handleResultFn,
	}

	// we don't want more pending requests than the number of peers we can query
	requestsEvents := concurrency
	if len(closestPeers) < concurrency {
		requestsEvents = len(closestPeers)
	}
	for i := 0; i < requestsEvents; i++ {
		// add concurrency requests to the event queue
		query.eventQueue.Enqueue(query.newRequest)
		span.AddEvent("Enqueued SimpleQuery.newRequest. Queue size: " + strconv.Itoa(int(query.eventQueue.Size())))
	}
	query.inflightRequests = requestsEvents

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
		close(q.resultsChan)
		return errors.New("query cancelled")
	default:
	}
	return nil
}

func (q *SimpleQuery) newRequest(ctx context.Context) {
	ctx, span := internal.StartSpan(ctx, "SimpleQuery.newRequest")
	defer span.End()

	if err := q.checkIfDone(); err != nil {
		span.RecordError(err)
		q.inflightRequests--
		return
	}

	peerid := popClosestQueued(q.peerlist)
	if peerid == "" {
		// TODO: handle this case
		span.AddEvent("all peers queried")
		q.inflightRequests--
		return
	}
	span.AddEvent("peer selected: " + peerid.String())

	// start new go routine to send request to peer
	go q.sendRequest(ctx, peerid)

	// add timeout to scheduler
	events.ScheduleAction(ctx, &q.eventPlanner, q.timeout, func(ctx context.Context) {
		q.requestError(ctx, peerid, errors.New("request timeout ("+q.timeout.String()+")"))
	})
}

// sendRequest
// note that this function is called in a separate go routine
func (q *SimpleQuery) sendRequest(ctx context.Context, p peer.ID) {
	ctx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()

	ctx, span := internal.StartSpan(ctx, "SimpleQuery.sendRequest")
	defer span.End()

	if err := q.msgEndpoint.DialPeer(ctx, p); err != nil {
		span.AddEvent("peer dial failed")
		q.eventQueue.Enqueue(func(ctx context.Context) {
			q.requestError(ctx, p, err)
		})
		return
	}

	resp, err := q.msgEndpoint.SendRequest(ctx, p, q.message, q.proto)
	if err != nil {
		span.AddEvent("request failed")
		q.eventQueue.Enqueue(func(ctx context.Context) {
			q.requestError(ctx, p, err)
		})
		return
	}
	span.AddEvent("got a response")

	q.eventQueue.Enqueue(func(ctx context.Context) {
		q.handleResponse(ctx, p, resp)
	})
	span.AddEvent("Enqueued SimpleQuery.handleResponse. Queue size: " + strconv.Itoa(int(q.eventQueue.Size())))
}

func (q *SimpleQuery) handleResponse(ctx context.Context, p peer.ID, resp *pb.Message) {
	ctx, span := internal.StartSpan(ctx, "SimpleQuery.handleResponse",
		trace.WithAttributes(attribute.String("Target", q.kadid.Hex()), attribute.String("From Peer", p.String())))
	defer span.End()

	if err := q.checkIfDone(); err != nil {
		span.RecordError(err)
		return
	}

	closerPeers := resp.GetCloserPeers()
	if len(closerPeers) > 0 {
		// consider that remote peer is behaving correctly if it returns
		// at least 1 peer
		q.rt.AddPeer(ctx, p)
	}

	updatePeerStatusInPeerlist(q.peerlist, p, queried)

	newPeers := network.ParsePeers(ctx, closerPeers)
	newPeerIds := make([]peer.ID, 0, len(newPeers))

	for _, ai := range newPeers {
		if key.Compare(key.PeerKadID(ai.ID), q.self) == 0 {
			// don't add self to queries or routing table
			span.AddEvent("remote peer provided self as closer peer")
			continue
		}
		newPeerIds = append(newPeerIds, ai.ID)

		q.msgEndpoint.MaybeAddToPeerstore(ai, QueuedPeersPeerstoreTTL)
	}

	addToPeerlist(q.peerlist, newPeerIds)

	var stop bool
	q.state = q.handleResultFn(ctx, q.state, resp, q.resultsChan)
	if stop {
		// query is done, don't send any more requests
		span.AddEvent("query success")
		q.done = true
		return
	}

	// we always want to have the maximal number of requests in flight
	newRequestsToSend := 1 + q.concurrency - q.inflightRequests
	if q.peerlist.queuedCount < newRequestsToSend {
		newRequestsToSend = q.peerlist.queuedCount
	}

	for i := 0; i < newRequestsToSend; i++ {
		// add new pending request(s) for this query to eventqueue
		q.eventQueue.Enqueue(q.newRequest)

	}
	span.AddEvent("Enqueued " + strconv.Itoa(newRequestsToSend) +
		"SimpleQuery.newRequest. Queue size: " +
		strconv.Itoa(int(q.eventQueue.Size())))
}

func (q *SimpleQuery) requestError(ctx context.Context, peerid peer.ID, err error) {
	ctx, span := internal.StartSpan(ctx, "SimpleQuery.requestError",
		trace.WithAttributes(attribute.String("PeerID", peerid.String()),
			attribute.String("Error", err.Error())))
	defer span.End()

	if q.ctx.Err() == nil {
		// remove peer from routing table unless context was cancelled
		q.rt.RemovePeer(ctx, key.PeerKadID(peerid))
	}

	if err := q.checkIfDone(); err != nil {
		span.RecordError(err)
		return
	}

	updatePeerStatusInPeerlist(q.peerlist, peerid, unreachable)

	// add pending request for this query to eventqueue
	q.eventQueue.Enqueue(q.newRequest)
}
