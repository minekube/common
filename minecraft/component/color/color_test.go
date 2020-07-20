package color

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHex(t *testing.T) {
	c, err := ParseHexString("#000000")
	require.NoError(t, err)
	require.Equal(t, Color{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}, c)

	c2, err := ParseHexString("#ffaa00")
	require.NoError(t, err)
	require.Equal(t, Color{
		R: 0xff,
		G: 0xaa,
		B: 0x00,
		A: 0xff,
	}, c2)

	require.Equal(t, "#ffaa00", c2.HexString())

	require.Equal(t, Color{
		R: 0xff,
		G: 0xaa,
		B: 0x00,
		A: 0xff,
	}, ParseHex(0xffaa00))
}
