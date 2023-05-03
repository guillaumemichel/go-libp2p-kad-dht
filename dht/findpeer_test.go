package dht

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/routing"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	libp2prouting "github.com/libp2p/go-libp2p/core/routing"
	"github.com/stretchr/testify/require"
)

func newDhtHost(ctx context.Context) *IpfsDHT {
	var idht *IpfsDHT

	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	// create new libp2p node using the keypair and set the DHT as the routing
	// system
	_, err = libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		libp2p.Routing(func(h host.Host) (libp2prouting.PeerRouting, error) {
			idht = NewDHT(ctx, h)
			return idht, err
		}),
	)
	if err != nil {
		panic(err)
	}
	return idht
}

func TestFindPeerQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dht0 := newDhtHost(ctx)

	// rtng (reference peer) only knows about rpeer0
	// rpeer0 knows about rpeer1
	// rtng can ask rpeer0 about rpeer1

	nRemotePeers := 2
	remotePeers := make([]*IpfsDHT, nRemotePeers)
	for i := 0; i < nRemotePeers; i++ {
		remotePeers[i] = newDhtHost(ctx)
	}
	dht0.Host.Peerstore().AddAddrs(remotePeers[0].Host.ID(),
		remotePeers[0].Host.Addrs(), routing.PEERSTORE_ENTRY_TTL)
	remotePeers[0].Host.Peerstore().AddAddrs(remotePeers[1].Host.ID(),
		remotePeers[1].Host.Addrs(), routing.PEERSTORE_ENTRY_TTL)

	// find a peer
	res, err := dht0.FindPeer(ctx, remotePeers[1].Host.ID())
	require.NoError(t, err)
	require.NotNil(t, res)
}
