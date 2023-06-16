package ipfskadv1

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-kad-dht/key"
	"github.com/libp2p/go-libp2p-kad-dht/network/address"
	laddr "github.com/libp2p/go-libp2p-kad-dht/network/address/libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/network/endpoint"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ErrNoValidAddresses = errors.New("no valid addresses")
)

func (msg *Message) Target() key.KadKey {
	p, err := peer.IDFromBytes(msg.GetKey())
	if err != nil {
		return key.KadKey{}
	}
	return key.PeerKadID(p)
}

func (msg *Message) CloserNodes() []address.NetworkAddress {
	closerPeers := msg.GetCloserPeers()
	if closerPeers == nil {
		return []address.NetworkAddress{}
	}
	return ParsePeers(closerPeers)
}

func PBPeerToPeerInfo(pbp *Message_Peer) (peer.AddrInfo, error) {
	addrs := make([]multiaddr.Multiaddr, 0, len(pbp.Addrs))
	for _, a := range pbp.Addrs {
		addr, err := multiaddr.NewMultiaddrBytes(a)
		if err == nil {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return peer.AddrInfo{}, ErrNoValidAddresses
	}

	return peer.AddrInfo{
		ID:    peer.ID(pbp.Id),
		Addrs: addrs,
	}, nil
}

func ParsePeers(pbps []*Message_Peer) []address.NetworkAddress {

	peers := make([]address.NetworkAddress, 0, len(pbps))
	for _, p := range pbps {
		pi, err := PBPeerToPeerInfo(p)
		if err == nil {
			peers = append(peers, laddr.Libp2pAddr(pi))
		}
	}
	return peers
}

func FindPeerRequest(p peer.ID) *Message {
	marshalledPeerid, _ := p.MarshalBinary()
	return &Message{
		Type: Message_FIND_NODE,
		//ClusterLevelRaw: 1,
		Key: marshalledPeerid,
	}
}

func FindPeerResponse(p peer.ID, peers []address.NodeID, e endpoint.Endpoint) *Message {
	marshalledPeerid, _ := p.MarshalBinary()
	return &Message{
		Type:        Message_FIND_NODE,
		Key:         marshalledPeerid,
		CloserPeers: NodeIDsToPbPeers(peers, e),
	}
}

func NodeIDsToPbPeers(peers []address.NodeID, e endpoint.Endpoint) []*Message_Peer {
	pbPeers := make([]*Message_Peer, 0, len(peers))
	for _, n := range peers {
		p := n.(peer.ID)

		na, err := e.NetworkAddress(n)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// convert NetworkAddress to []multiaddr.Multiaddr
		addrs := na.(peer.AddrInfo).Addrs
		pbAddrs := make([][]byte, len(addrs))
		// convert multiaddresses to bytes
		for i, a := range addrs {
			pbAddrs[i] = a.Bytes()
		}

		pbPeers = append(pbPeers, &Message_Peer{
			Id:         []byte(p),
			Addrs:      pbAddrs,
			Connection: Message_ConnectionType(e.Connectedness(n)),
		})
	}
	return pbPeers
}

func PeeridsToPbPeers(peers []peer.ID, h host.Host) []*Message_Peer {

	pbPeers := make([]*Message_Peer, 0, len(peers))

	for _, p := range peers {
		addrs := h.Peerstore().Addrs(p)
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
			Id:         []byte(p),
			Addrs:      pbAddrs,
			Connection: Message_ConnectionType(h.Network().Connectedness(p)),
		})
	}
	return pbPeers
}
