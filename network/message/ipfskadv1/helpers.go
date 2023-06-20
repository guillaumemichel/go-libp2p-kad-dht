package ipfskadv1

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/addrinfo"
	"github.com/libp2p/go-libp2p-kad-dht/network/address/peerid"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ErrNoValidAddresses = errors.New("no valid addresses")
)

func FindPeerRequest(p peerid.PeerID) *Message {
	marshalledPeerid, _ := p.MarshalBinary()
	return &Message{
		Type: Message_FIND_NODE,
		Key:  marshalledPeerid,
	}
}

func FindPeerResponse(peers []address.NodeID, e endpoint.NetworkedEndpoint) *Message {
	return &Message{
		Type:        Message_FIND_NODE,
		CloserPeers: NodeIDsToPbPeers(peers, e),
	}
}

func (msg *Message) Target() key.KadKey {
	p, err := peer.IDFromBytes(msg.GetKey())
	if err != nil {
		return nil
	}
	return peerid.PeerID{ID: p}.Key()
}

func (msg *Message) CloserNodes() []address.NetworkAddress {
	closerPeers := msg.GetCloserPeers()
	if closerPeers == nil {
		return []address.NetworkAddress{}
	}
	return ParsePeers(closerPeers)
}

func PBPeerToPeerInfo(pbp *Message_Peer) (addrinfo.AddrInfo, error) {
	addrs := make([]multiaddr.Multiaddr, 0, len(pbp.Addrs))
	for _, a := range pbp.Addrs {
		addr, err := multiaddr.NewMultiaddrBytes(a)
		if err == nil {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return addrinfo.AddrInfo{}, ErrNoValidAddresses
	}

	return addrinfo.AddrInfo{
		AddrInfo: peer.AddrInfo{
			ID:    peer.ID(pbp.Id),
			Addrs: addrs,
		},
	}, nil
}

func ParsePeers(pbps []*Message_Peer) []address.NetworkAddress {

	peers := make([]address.NetworkAddress, 0, len(pbps))
	for _, p := range pbps {
		pi, err := PBPeerToPeerInfo(p)
		if err == nil {
			peers = append(peers, pi)
		}
	}
	return peers
}

func NodeIDsToPbPeers(peers []address.NodeID, e endpoint.NetworkedEndpoint) []*Message_Peer {
	if len(peers) == 0 || e == nil {
		return nil
	}

	pbPeers := make([]*Message_Peer, 0, len(peers))
	for _, n := range peers {
		p := n.(peerid.PeerID)

		na, err := e.NetworkAddress(n)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// convert NetworkAddress to []multiaddr.Multiaddr
		addrs := na.(addrinfo.AddrInfo).Addrs
		pbAddrs := make([][]byte, len(addrs))
		// convert multiaddresses to bytes
		for i, a := range addrs {
			pbAddrs[i] = a.Bytes()
		}

		pbPeers = append(pbPeers, &Message_Peer{
			Id:         []byte(p.ID),
			Addrs:      pbAddrs,
			Connection: Message_ConnectionType(e.Connectedness(n)),
		})
	}
	return pbPeers
}

func PeeridsToPbPeers(peers []peerid.PeerID, h host.Host) []*Message_Peer {

	pbPeers := make([]*Message_Peer, 0, len(peers))

	for _, p := range peers {
		addrs := h.Peerstore().Addrs(p.ID)
		if len(addrs) == 0 {
			// if no addresses, don't send peer
			continue
		}

		pbAddrs := make([][]byte, len(addrs))
		// convert multiaddresses to bytes
		for i, a := range addrs {
			pbAddrs[i] = a.Bytes()
		}
		pbPeers = append(pbPeers, &Message_Peer{
			Id:         []byte(p.ID),
			Addrs:      pbAddrs,
			Connection: Message_ConnectionType(h.Network().Connectedness(p.ID)),
		})
	}
	return pbPeers
}
