package main

// test program. we can use the go one.
import (
	"context"
	"fmt"
	"time"

	// varint is here

	"github.com/libp2p/go-libp2p/p2p/net/connmgr"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/dht"
	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

func main() {
	ctx := context.Background()

	var dht0, dht1 *dht.IpfsDHT
	host0, err := libp2pHost(ctx, "10000", dht0)
	if err != nil {
		panic(err)
	}
	host1, err := libp2pHost(ctx, "10001", dht1)
	if err != nil {
		panic(err)
	}

	host0.Peerstore().AddAddrs(host1.ID(), host1.Addrs(), peerstore.TempAddrTTL)
	host1.Peerstore().AddAddrs(host0.ID(), host0.Addrs(), peerstore.TempAddrTTL)

	if err := host0.Connect(ctx, host0.Peerstore().PeerInfo(host1.ID())); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connected")

	s, err := host0.NewStream(ctx, host1.ID(), protocol.ProtocolDHT)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	message := []byte("hello world")
	_, err = s.Write(message)
	if err != nil {
		fmt.Println("error host0", err)
	}

	time.Sleep(1 * time.Second)

	s.Write([]byte("hello world"))

	time.Sleep(5 * time.Second)

	host0.Close()
	host1.Close()
}

func libp2pHost(ctx context.Context, port string, idht *dht.IpfsDHT) (host.Host, error) {
	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}
	h2, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+port,         // regular tcp connections
			"/ip4/0.0.0.0/udp/"+port+"/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht = dht.NewDHT(ctx, h)
			return idht, err
		}),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}
	return h2, nil
}
