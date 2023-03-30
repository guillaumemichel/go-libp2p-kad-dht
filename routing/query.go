package routing

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
)

type query struct {
	// query peer set
	// {dist, conn_type, peer.Info, queried?}

	ctx context.Context

	kadId   hash.KadKey
	results chan queryResult

	routing *DhtRouting
}

type queryResult struct {
	// errors
	// peers
	// values
}

type queryManager struct {
	limit          int
	ongoingQueries []*query
	queuedQueries  []*query

	routing *DhtRouting
}

func (r *DhtRouting) newQueryManager(limit int) *queryManager {
	return &queryManager{
		limit:          limit,
		ongoingQueries: []*query{},
		queuedQueries:  []*query{},
		routing:        r,
	}
}

// can return either a peer or a value (or both).
// peer is peerid + multiaddrs or peerid only (peer.AddrInfo)
// value is a byte array (or string)

func (qm *queryManager) newQuery(ctx context.Context, kadId hash.KadKey) *query {
	return &query{
		ctx:     ctx,
		kadId:   kadId,
		routing: qm.routing,
		results: make(chan queryResult),
	}
}

func (qm *queryManager) Query(ctx context.Context, kadId hash.KadKey, msg *pb.DhtMessage, stopCond func() bool) chan queryResult {
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

	//sollicitedPeers := r.rt.NearestPeers(q.kadId, r.concurrency)

}
