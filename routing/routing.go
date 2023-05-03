package routing

import (
	"context"

	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	QUERY_LIMIT = 30
)

type PeerRouting interface {
	FindPeer(context.Context, peer.ID) (peer.AddrInfo, error)
}

type ContentRouting interface {
}

type DhtRouting struct {
	ctx         context.Context
	me          *dhtnet.MessageEndpoint
	rt          simplert.RoutingTable
	concurrency int
	qManager    *queryManager
}

func NewDhtRouting(ctx context.Context, me *dhtnet.MessageEndpoint, rt simplert.RoutingTable, concurrency int, queryLimit int) *DhtRouting {
	r := &DhtRouting{
		ctx:         ctx,
		me:          me,
		rt:          rt,
		concurrency: concurrency,
	}
	r.qManager = r.newQueryManager(ctx, queryLimit)
	return r
}
