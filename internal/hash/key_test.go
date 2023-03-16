package hash

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	mhreg "github.com/multiformats/go-multihash/core"
	"github.com/stretchr/testify/assert"
)

func TestPeerKadID(t *testing.T) {
	peerid := peer.ID("12D3KooWFHYCmTCexEziKBVFEVT4FAVPBUYJKEBUrdgu9RYoVE1T")
	kadid := PeerKadID(peerid)

	// get sha256 hasher
	hasher, err := mhreg.GetHasher(HasherID)
	assert.NoError(t, err)
	hasher.Write([]byte(peerid))
	expectedKadid := hasher.Sum(nil)

	assert.Equal(t, kadid[:], expectedKadid)
}
