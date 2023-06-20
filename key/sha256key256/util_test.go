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
