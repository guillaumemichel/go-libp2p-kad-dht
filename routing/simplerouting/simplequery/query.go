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
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	message "github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable"
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

type HandleResultFn func(context.Context, []interface{}, message.MinKadResponseMessage, chan interface{}) []interface{}

type SimpleQuery struct {
	ctx         context.Context
	done        bool
	kadid       key.KadKey
	req         message.MinKadRequestMessage
	resp        message.MinKadResponseMessage
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
func NewSimpleQuery(ctx context.Context, kadid key.KadKey, req message.MinKadRequestMessage,
	resp message.MinKadResponseMessage, concurrency int, timeout time.Duration,
	proto protocol.ID, msgEndpoint endpoint.Endpoint, rt routingtable.RoutingTable,
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
		req:              req,
		resp:             resp,
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

	id := popClosestQueued(q.peerlist)
	if id == nil || id.String() == "" {
		// TODO: handle this case
		span.AddEvent("all peers queried")
		q.inflightRequests--
		return
	}
	span.AddEvent("peer selected: " + id.String())

	// start new go routine to send request to peer
	go q.sendRequest(ctx, id)

	// add timeout to scheduler
	events.ScheduleAction(ctx, &q.eventPlanner, q.timeout, func(ctx context.Context) {
		q.requestError(ctx, id, errors.New("request timeout ("+q.timeout.String()+")"))
	})
}

// sendRequest
// note that this function is called in a separate go routine
func (q *SimpleQuery) sendRequest(ctx context.Context, id address.NodeID) {
	ctx, cancel := context.WithTimeout(ctx, q.timeout)
	defer cancel()

	ctx, span := internal.StartSpan(ctx, "SimpleQuery.sendRequest")
	defer span.End()

	if err := q.msgEndpoint.DialPeer(ctx, id); err != nil {
		span.AddEvent("peer dial failed")
		q.eventQueue.Enqueue(func(ctx context.Context) {
			q.requestError(ctx, id, err)
		})
		return
	}

	err := q.msgEndpoint.SendRequest(ctx, id, q.req, q.resp, q.proto)
	if err != nil {
		span.AddEvent("request failed")
		q.eventQueue.Enqueue(func(ctx context.Context) {
			q.requestError(ctx, id, err)
		})
		return
	}
	span.AddEvent("got a response")

	q.eventQueue.Enqueue(func(ctx context.Context) {
		q.handleResponse(ctx, id, q.resp)
	})
	span.AddEvent("Enqueued SimpleQuery.handleResponse. Queue size: " + strconv.Itoa(int(q.eventQueue.Size())))
}

func (q *SimpleQuery) handleResponse(ctx context.Context, id address.NodeID, resp message.MinKadResponseMessage) {
	ctx, span := internal.StartSpan(ctx, "SimpleQuery.handleResponse",
		trace.WithAttributes(attribute.String("Target", q.kadid.Hex()), attribute.String("From Peer", id.String())))
	defer span.End()

	if err := q.checkIfDone(); err != nil {
		span.RecordError(err)
		return
	}

	closerPeers := resp.CloserNodes()
	if len(closerPeers) > 0 {
		// consider that remote peer is behaving correctly if it returns
		// at least 1 peer
		q.rt.AddPeer(ctx, id)
	}

	updatePeerStatusInPeerlist(q.peerlist, id, queried)

	//newPeers := message.ParsePeers(ctx, closerPeers)
	newPeerIds := make([]address.NodeID, 0, len(closerPeers))

	for _, na := range closerPeers {
		id := na.NodeID()
		if key.Compare(address.KadID(id), q.msgEndpoint.KadID()) == 0 {
			// don't add self to queries or routing table
			span.AddEvent("remote peer provided self as closer peer")
			continue
		}
		newPeerIds = append(newPeerIds, id)

		q.msgEndpoint.MaybeAddToPeerstore(na, QueuedPeersPeerstoreTTL)
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
		" SimpleQuery.newRequest. Queue size: " +
		strconv.Itoa(int(q.eventQueue.Size())))

}

func (q *SimpleQuery) requestError(ctx context.Context, id address.NodeID, err error) {
	ctx, span := internal.StartSpan(ctx, "SimpleQuery.requestError",
		trace.WithAttributes(attribute.String("PeerID", id.String()),
			attribute.String("Error", err.Error())))
	defer span.End()

	if q.ctx.Err() == nil {
		// remove peer from routing table unless context was cancelled
		q.rt.RemovePeer(ctx, address.KadID(id))
	}

	if err := q.checkIfDone(); err != nil {
		span.RecordError(err)
		return
	}

	updatePeerStatusInPeerlist(q.peerlist, id, unreachable)

	// add pending request for this query to eventqueue
	q.eventQueue.Enqueue(q.newRequest)
}
