package simplerouting

import (
	"context"

	"github.com/libp2p/go-libp2p-kad-dht/internal"
	"github.com/libp2p/go-libp2p-kad-dht/internal/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/message"
	"github.com/libp2p/go-libp2p-kad-dht/network/message/ipfskadv1"
	sq "github.com/libp2p/go-libp2p-kad-dht/routing/simplerouting/simplequery"
	libp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// when trying to write results to the results channel, we need to make sure
// that we don't block forever. select on ctx
type HandleResultsFn func(context.Context, sq.SimpleQuery)

// FindPeer searches for a peer with given ID. This is a blocking call.
func (r *SimpleRouting) FindPeer(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	ctx, span := internal.StartSpan(ctx, "SimpleRouting.FindPeer",
		trace.WithAttributes(attribute.String("PeerID", p.String())))
	defer span.End()

	if err := p.Validate(); err != nil {
		return peer.AddrInfo{}, err
	}

	// Check if were already connected to them
	targetConnectedness := r.msgEndpoint.Connectedness(p)
	if targetConnectedness == libp2pnet.Connected {
		span.AddEvent("Already connected")
		return r.msgEndpoint.PeerInfo(p), nil
	}

	kadid := key.PeerKadID(p)
	req := ipfskadv1.FindPeerRequest(p)
	resp := &ipfskadv1.Message{}

	resultsChan := make(chan interface{}) // peer.AddrInfo
	handleResultsFn := getFindPeerHandleResultsFn(p)

	// this serve to cancel the query (dependant on ctx) once we return a result
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create the query and add appropriate events to the event queue
	sq.NewSimpleQuery(ctx, kadid, req, resp, r.queryConcurrency, r.queryTimeout,
		r.protocolID, r.msgEndpoint, r.rt, r.eventQueue, r.eventPlanner,
		resultsChan, handleResultsFn)

	// only one dial runs at a time to ensure sequentiality
	dialRunning := false
	newDialRequired := false
	dialChan := make(chan bool)
	currAddresses := r.msgEndpoint.PeerInfo(p).Addrs

	// this function is called once the async dial finishes, reporting its result
	dialReportFn := func(ctx context.Context, success bool) {
		select {
		case <-ctx.Done():
		case dialChan <- success:
		}
	}

	// make sure we can still connect to the peer at the address we have
	if targetConnectedness == libp2pnet.CanConnect {
		span.AddEvent("Already in peerstore: can connect")
		dialRunning = true
		// spawns a new goroutine to dial the peer and report the result
		r.msgEndpoint.AsyncDialAndReport(ctx, p, dialReportFn)
	}

	select {
	case <-ctx.Done():
		// query was cancelled
		return peer.AddrInfo{}, ctx.Err()

	case res := <-resultsChan:
		// we got a result from the query, we need to check if the address is
		// valid
		ai, ok := res.(peer.AddrInfo)
		if !ok {
			span.AddEvent("Unexpected result type")
		} else {
			var newAddr bool
			newAddr, currAddresses = containsNewAddresses(ai.Addrs, currAddresses)
			if newAddr {
				// if we found a new address, we need to dial the peer to make
				// sure the new address is valid
				if dialRunning {
					newDialRequired = true
				} else {
					dialRunning = true
					// spawns a new goroutine to dial the peer and report the result
					r.msgEndpoint.AsyncDialAndReport(ctx, p, dialReportFn)
				}
			}
		}
	case success := <-dialChan:
		if success {
			// if we could dial the peer, return its address info
			return r.msgEndpoint.PeerInfo(p), nil
		}
		if newDialRequired {
			newDialRequired = false
			// spawns a new goroutine to dial the peer and report the result
			r.msgEndpoint.AsyncDialAndReport(ctx, p, dialReportFn)
		} else {
			dialRunning = false
		}
	}
	return r.msgEndpoint.PeerInfo(p), nil
}

// containsNewAddresses returns true if newAddrs contains addresses that are not
// in oldAddrs. It also returns the union of both address slices.
func containsNewAddresses(newAddrs, oldAddrs []multiaddr.Multiaddr) (bool, []multiaddr.Multiaddr) {
	bNewAddr := false
	for _, n := range newAddrs {
		found := false
		for _, o := range oldAddrs {
			if n.Equal(o) {
				found = true
				break
			}
		}
		if !found {
			bNewAddr = true
			oldAddrs = append(oldAddrs, n)
		}
	}
	return bNewAddr, oldAddrs
}

// getFindPeerHandleResultsFn returns a HandleResultsFn that checks if any
// peer.ID of the result matches the peer.ID we are looking for. If one does,
// it writes the result to the resultsChan and returns nil
func getFindPeerHandleResultsFn(p peer.ID) sq.HandleResultFn {
	return func(ctx context.Context, i []interface{}, m message.MinKadResponseMessage,
		resultsChan chan interface{}) []interface{} {

		ctx, span := internal.StartSpan(ctx, "SimpleRouting.getFindPeerHandleResultsFn")
		defer span.End()

		for _, na := range m.CloserNodes() {
			if na.NodeID().String() == p.String() {
				// we found the peer we were looking for

				// convert NodeID to PeerID as we need to return a PeerID to the caller
				peerid := na.NodeID().(peer.ID)
				span.AddEvent("Found peer", trace.WithAttributes(attribute.String("PeerID", peerid.String())))

				select {
				case <-ctx.Done():
					return nil
				case resultsChan <- peerid:
				}

				break
			}
		}
		return nil
	}
}
