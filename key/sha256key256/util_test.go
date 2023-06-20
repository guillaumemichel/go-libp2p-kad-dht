package sha256key256

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZeroKey(t *testing.T) {
	zero := ZeroKey()
	require.Equal(t, 2*Keysize, len(zero.String()))

	xored, _ := zero.Xor(zero)
	require.Equal(t, zero, xored)
}

func TestPeerKadID(t *testing.T) {
	str := "hello world"
	// generated with: $ echo -n "hello world" | sha256sum
	digest := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	kid := StringKadID(str)
	require.Equal(t, digest, kid.String())
}
