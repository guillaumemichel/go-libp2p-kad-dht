package endpoint

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1/pb"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type DialReportFn func(context.Context, bool)

type Endpoint interface {
	AsyncDialAndReport(ctx context.Context, p peer.ID, reportFn DialReportFn)
	DialPeer(ctx context.Context, p peer.ID) error
	MaybeAddToPeerstore(ai peer.AddrInfo, ttl time.Duration)
	SendRequest(ctx context.Context, p peer.ID, req *pb.Message, proto protocol.ID) (*pb.Message, error)

	// Peerstore functions
	KadID() key.KadKey
	Connectedness(p peer.ID) network.Connectedness
	PeerInfo(p peer.ID) peer.AddrInfo
}
