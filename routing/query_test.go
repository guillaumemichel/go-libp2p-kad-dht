package routing

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/dht/protocol"
	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	dhtnet "github.com/libp2p/go-libp2p-kad-dht/network"
	"github.com/libp2p/go-libp2p-kad-dht/network/pb"
	"github.com/libp2p/go-libp2p-kad-dht/routingtable/simplert"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/stretchr/testify/require"
)

func zeroBytes(n int) []byte {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = 0
	}
	return bytes
}

var (
	keys = []hash.KadKey{
		hash.KadKey(zeroBytes(hash.Keysize)),                // 0000 0000 ... 0000 0000
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x1)), // 0000 0000 ... 0000 0001
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x2)), // 0000 0000 ... 0000 0010
		hash.KadKey(append(zeroBytes(hash.Keysize-1), 0x3)), // 0000 0000 ... 0000 0011
	}
)

func TestQpeersetInsert(t *testing.T) {
	q := newQpeerset(keys[0])
	p0 := &qpeer{dist: keys[0]}
	p1 := &qpeer{dist: keys[1]}
	p2 := &qpeer{dist: keys[2]}
	p3 := &qpeer{dist: keys[3]}

	require.Equal(t, 0, q.size)
	q.insert(p1)
	require.Equal(t, 1, q.size)
	q.insert(p1)
	require.Equal(t, 1, q.size)
	q.insert(p0)
	require.Equal(t, 2, q.size)
	require.Equal(t, p0.dist, q.head.dist)
	require.Equal(t, p1.dist, q.head.next.dist)

	q.insert(p3)
	q.insert(p2)
	require.Equal(t, 4, q.size)
	require.Equal(t, p2.dist, q.head.next.next.dist)

	dummyAddrInfo := peer.AddrInfo{ID: "dummy"}

	qpeer := q.AddPeer(dummyAddrInfo)
	require.NotNil(t, qpeer)
	require.Equal(t, 5, q.size)
	qpeer = q.AddPeer(dummyAddrInfo)
	require.Nil(t, qpeer)
	require.Equal(t, 5, q.size)

}

func TestClosestHeard(t *testing.T) {
	q := newQpeerset(keys[0])
	p0 := &qpeer{dist: keys[0]}
	p1 := &qpeer{dist: keys[1]}
	p2 := &qpeer{dist: keys[2]}
	p3 := &qpeer{dist: keys[3]}

	q.insert(p1)
	q.insert(p0)
	q.insert(p3)
	q.insert(p2)

	require.Equal(t, []*qpeer{p0, p1, p2}, q.ClosestHeard(3))
}

func newTestDhtRouting(ctx context.Context) *DhtRouting {

	// start a libp2p node with default settings
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	// create message endpoint
	me := dhtnet.NewMessageEndpoint(node)
	// create a dht routing

	// create a dht routing table
	rt := simplert.NewDhtRoutingTable(hash.PeerKadID(node.ID()), 20)

	return NewDhtRouting(ctx, me, rt, 1, 1)
}

func basicStreamHandler(s network.Stream) {
	// create a protobuf reader and writer
	r := pbio.NewDelimitedReader(s, network.MessageSizeMax)
	w := pbio.NewDelimitedWriter(s)

	for {
		req := &pb.DhtMessage{}
		// read an empty message from the stream
		err := r.ReadMsg(req)
		if err != nil {
			if err == io.EOF {
				// stream EOF, all done
				return
			}
			fmt.Println(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
		resp := &pb.DhtMessage{}
		// write an empty response to the stream
		err = w.WriteMsg(resp)
		if err != nil {
			return
		}
	}
}

func newTestRemotePeer(ctx context.Context) host.Host {
	// start a libp2p node with default settings
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(protocol.ProtocolDHT, basicStreamHandler)
	return node
}

func TestQueryManagerDisconnected(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rtng := newTestDhtRouting(ctx)

	dummyMsg := &pb.DhtMessage{}

	require.Equal(t, 0, rtng.qManager.currentlyOngoing)

	res0 := rtng.qManager.Query(ctx, keys[0], dummyMsg)
	require.Equal(t, 1, rtng.qManager.currentlyOngoing)

	count := 1

	for count > 0 {
		select {
		case <-ctx.Done():
			return
		case <-res0:
			count--
		}
	}
	time.Sleep(1 * time.Millisecond) // give time to query manager to update the counter
	require.Equal(t, 0, rtng.qManager.currentlyOngoing)
}

func TestQueryManager(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rtng := newTestDhtRouting(ctx)

	nRemotePeers := 10
	remotePeers := make([]host.Host, nRemotePeers)
	for i := 0; i < nRemotePeers; i++ {
		remotePeers[i] = newTestRemotePeer(ctx)

		rtng.me.Host.Peerstore().AddAddrs(remotePeers[i].ID(), remotePeers[i].Addrs(), PEERSTORE_ENTRY_TTL)
		rtng.rt.AddPeer(peer.AddrInfo{ID: remotePeers[i].ID(), Addrs: remotePeers[i].Addrs()})
	}

	dummyMsg := &pb.DhtMessage{}
	res0 := rtng.qManager.Query(ctx, keys[0], dummyMsg)
	res1 := rtng.qManager.Query(ctx, keys[1], dummyMsg)

	shortCtx, shortCancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer shortCancel()

	res2 := rtng.qManager.Query(shortCtx, keys[2], dummyMsg)
	count := 3

	for count > 0 {
		select {
		case <-ctx.Done():
			return
		case r := <-res0:
			require.Equal(t, ErrUnexpectedResponseType, r.err)
			count--
		case r := <-res1:
			require.Equal(t, ErrUnexpectedResponseType, r.err)
			count--
		case r := <-res2:
			require.Equal(t, ErrUnexpectedResponseType, r.err)
			count--
		}
	}
}

func TestQueryManagerCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	rtng := newTestDhtRouting(ctx)

	nRemotePeers := 10
	remotePeers := make([]host.Host, nRemotePeers)
	for i := 0; i < nRemotePeers; i++ {
		remotePeers[i] = newTestRemotePeer(ctx)

		rtng.me.Host.Peerstore().AddAddrs(remotePeers[i].ID(), remotePeers[i].Addrs(), PEERSTORE_ENTRY_TTL)
		rtng.rt.AddPeer(peer.AddrInfo{ID: remotePeers[i].ID(), Addrs: remotePeers[i].Addrs()})
	}

	dummyMsg := &pb.DhtMessage{}
	res0 := rtng.qManager.Query(ctx, keys[0], dummyMsg)
	res1 := rtng.qManager.Query(ctx, keys[1], dummyMsg)

	cancel()

	time.Sleep(10 * time.Millisecond) // give time to query manager to update the counter

	require.Equal(t, 0, len(res0))
	require.Equal(t, 0, len(res1))
}

func TestHandleFindPeerResponse(t *testing.T) {
	// TODO
}
func TestHandlePeers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rtng := newTestDhtRouting(ctx)

	nRemotePeers := 10
	remotePeers := make([]host.Host, nRemotePeers)
	for i := 0; i < nRemotePeers; i++ {
		remotePeers[i] = newTestRemotePeer(ctx)

		rtng.me.Host.Peerstore().AddAddrs(remotePeers[i].ID(), remotePeers[i].Addrs(), PEERSTORE_ENTRY_TTL)
		rtng.rt.AddPeer(peer.AddrInfo{ID: remotePeers[i].ID(), Addrs: remotePeers[i].Addrs()})
	}

	peerAddrs := make([]peer.AddrInfo, nRemotePeers)
	for i := 0; i < nRemotePeers; i++ {
		peerAddrs[i] = peer.AddrInfo{ID: remotePeers[i].ID(), Addrs: remotePeers[i].Addrs()}
	}

	q0 := rtng.qManager.newQuery(ctx, keys[0])

	q0.handlePeers(peerAddrs)

	// TODO: add test with peers that have already been queried
}
