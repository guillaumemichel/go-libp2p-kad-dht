package routing

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p/core/host"
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
	host        host.Host
	rt          simplert.RoutingTable
	concurrency int
	qManager    *queryManager
}

func NewDhtRouting(host host.Host, rt simplert.RoutingTable, concurrency int) *DhtRouting {
	r := &DhtRouting{
		host:        host,
		rt:          rt,
		concurrency: concurrency,
	}
	r.qManager = r.newQueryManager(QUERY_LIMIT)
	return r
}
