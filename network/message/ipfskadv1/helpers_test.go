package ipfskadv1

import (
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func TestFindPeerRequest(t *testing.T) {
	p, err := peer.Decode("12D3KooWH6Qd1EW75ANiCtYfD51D6M7MiZwLQ4g8wEBpoEUnVYNz")
	require.NoError(t, err)

	pid := peerid.PeerID{ID: p}
	msg := FindPeerRequest(pid)

	require.Equal(t, msg.GetKey(), []byte(p))

	c, err := msg.Target().Compare(pid.Key())
	require.NoError(t, err)
	require.Equal(t, int8(0), c)

	require.Equal(t, 0, len(msg.CloserNodes()))
}

func TestFindPeerResponse(t *testing.T) {

}

func TestCornerCases(t *testing.T) {

}
