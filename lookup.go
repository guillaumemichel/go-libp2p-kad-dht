package dht

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hashing"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
)

// GetClosestPeers is a Kademlia 'node lookup' operation. Returns a channel of
// the K closest peers to the given key.
//
// If the context is canceled, this function will return the context error along
// with the closest K peers it has found so far.
func (dht *IpfsDHT) GetClosestPeers(ctx context.Context, key hashing.KadKey) ([]peer.ID, error) {

	lookupRes, err := dht.runLookupWithFollowup(ctx, key,
		func(ctx context.Context, p peer.ID) ([]*peer.AddrInfo, error) {
			// For DHT query command
			routing.PublishQueryEvent(ctx, &routing.QueryEvent{
				Type: routing.SendingQuery,
				ID:   p,
			})

			peers, err := dht.protoMessenger.GetClosestPeers(ctx, p, key)
			if err != nil {
				logger.Debugf("error getting closer peers: %s", err)
				return nil, err
			}

			// For DHT query command
			routing.PublishQueryEvent(ctx, &routing.QueryEvent{
				Type:      routing.PeerResponse,
				ID:        p,
				Responses: peers,
			})

			return peers, err
		},
		func() bool { return false },
	)

	if err != nil {
		return nil, err
	}

	if ctx.Err() == nil && lookupRes.completed {
		// refresh the cpl for this key as the query was successful
		dht.routingTable.ResetCplRefreshedAtForID(key[:], time.Now())
	}

	return lookupRes.peers, ctx.Err()
}
