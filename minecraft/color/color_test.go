package color

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHex(t *testing.T) {
	c, err := Hex("#000000")
	require.NoError(t, err)
	require.Equal(t, Color{}, c)

	exp := Color{R: 1, G: 0.6666666666666666, B: 0}

	c2, err := Hex("#ffaa00")
	require.NoError(t, err)
	require.Equal(t, exp, c2)

	require.Equal(t, "#ffaa00", c2.Hex())

	require.Equal(t, exp, HexInt(0xffaa00))
}

func TestNearest(t *testing.T) {
	nearGold := HexInt(0xffaa01)
	require.Equal(t, GoldColor, nearGold.NearestNamed().Color)
}
