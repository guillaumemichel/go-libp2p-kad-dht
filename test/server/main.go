package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dht"
	"github.com/libp2p/go-libp2p-kad-dht/test/util"
)

func main() {
	ctx := context.Background()
	var dht *dht.IpfsDHT
	host, dht, err := util.Libp2pHost(ctx, "8000", dht)
	if err != nil {
		panic(err)
	}
	fmt.Println(host.ID())

	time.Sleep(1 * time.Minute)
}
