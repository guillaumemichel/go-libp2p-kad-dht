package dht

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/provider"
	"github.com/libp2p/go-libp2p-kad-dht/routing"
	"github.com/libp2p/go-libp2p-kad-dht/server"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type IpfsDHT struct {
	host     host.Host
	Provider provider.Provider
	net      *network.DhtNetwork
	server   *server.DhtServer
	routing  *routing.DhtRouting

	ctx context.Context
}

func NewDHT(ctx context.Context, h host.Host) *IpfsDHT {
	fmt.Println("creating new dht")
	net := network.NewDhtNetwork(h)
	prov := provider.NewDhtProvider(net)
	serv := server.NewDhtServer(net)
	routing := routing.NewDhtRouting(h)
	dht := &IpfsDHT{
		host:     h,
		Provider: prov,
		server:   serv,
		routing:  routing,
		ctx:      ctx,
	}

	return dht
}

func (dht *IpfsDHT) FindPeer(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	return dht.routing.FindPeer(ctx, p)
}
