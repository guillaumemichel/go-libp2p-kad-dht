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
	me          *dhtnet.MessageEndpoint
	rt          simplert.RoutingTable
	concurrency int
	qManager    *queryManager
}

func NewDhtRouting(me *dhtnet.MessageEndpoint, rt simplert.RoutingTable, concurrency int) *DhtRouting {
	r := &DhtRouting{
		me:          me,
		rt:          rt,
		concurrency: concurrency,
	}
	r.qManager = r.newQueryManager(QUERY_LIMIT)
	return r
}
