package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZero(t *testing.T) {
	require.Zero(t, Zero[string]())
	require.Zero(t, Zero[int]())
	require.Zero(t, Zero[float32]())
	require.Zero(t, Zero[float64]())
	require.Zero(t, Zero[uint64]())
	require.Zero(t, Zero[[]any]())
}
