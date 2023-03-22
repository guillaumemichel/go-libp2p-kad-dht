package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dht"
	"github.com/libp2p/go-libp2p-kad-dht/test/util"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

func main() {
	ctx := context.Background()
	var sDht *dht.IpfsDHT
	serverHost, _, err := util.Libp2pHost(ctx, "8000", sDht)
	if err != nil {
		panic(err)
	}
	//sDht = dht.NewDHT(ctx, serverHost)

	var cDht *dht.IpfsDHT
	clientHost, cDht, err := util.Libp2pHost(ctx, "8001", cDht)
	if err != nil {
		panic(err)
	}

	clientHost.Peerstore().AddAddrs(serverHost.ID(), serverHost.Addrs(), peerstore.TempAddrTTL)

	cid := util.GenCid()
	err = cDht.Provider.Provide(serverHost.ID(), cid)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("client done")
}
