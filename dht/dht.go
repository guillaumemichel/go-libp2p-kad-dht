package dht

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/provider"
	"github.com/libp2p/go-libp2p-kad-dht/routing"
	rt "github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p-kad-dht/server"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type IpfsDHT struct {
	Host        host.Host
	Provider    provider.Provider
	server      *server.DhtServer
	routing     *routing.DhtRouting
	rt          rt.RoutingTable
	msgEndpoint *dhtnet.MessageEndpoint

	ctx context.Context
}

func NewDHT(ctx context.Context, h host.Host) *IpfsDHT {
	var rt rt.RoutingTable = rt.NewDhtRoutingTable(hash.PeerKadID(h.ID()), 20)
	msgEndpoint := dhtnet.NewMessageEndpoint(h)
	prov := provider.NewDhtProvider(msgEndpoint)
	serv := server.NewDhtServer(msgEndpoint, rt)
	routing := routing.NewDhtRouting(ctx, msgEndpoint, rt, 10, routing.QUERY_LIMIT)
	dht := &IpfsDHT{
		Host:        h,
		rt:          rt,
		Provider:    prov,
		server:      serv,
		routing:     routing,
		msgEndpoint: msgEndpoint,
		ctx:         ctx,
	}

	return dht
}

func (dht *IpfsDHT) FindPeer(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	return dht.routing.FindPeer(ctx, p)
}