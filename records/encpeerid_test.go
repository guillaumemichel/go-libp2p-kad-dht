package records

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/libp2p/go-libp2p-kad-dht/internal/hash"
	"github.com/stretchr/testify/require"
)

func TestNonceSize(t *testing.T) {
	// key is the hash of the given string
	key := hash.PeerKadID("QmT8JQZj")
	block, err := aes.NewCipher(key[:])
	require.NoError(t, err)
	aead, err := cipher.NewGCM(block)
	require.NoError(t, err)
	require.Equal(t, aead.NonceSize(), NonceSize)
}
