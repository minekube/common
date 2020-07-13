package text

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_size(t *testing.T) {
	require.Equal(t, len(formats), len(LegacyChars))
}

func Test_lastIndexFrom(t *testing.T) {
	require.Equal(t, 3, lastIndexFrom("abca", 'a', 4))
	require.Equal(t, 3, lastIndexFrom("abca", 'a', 5))
	require.Equal(t, 0, lastIndexFrom("abca", 'a', 1))
	require.Equal(t, -1, lastIndexFrom("abca", 'a', 0))
}

func Test_reverseSlice(t *testing.T) {
	a := []*Text{
		{Text: "1"},
		{Text: "2"},
		{Text: "3"},
	}
	e := []*Text{
		{Text: "3"},
		{Text: "2"},
		{Text: "1"},
	}
	reverseSlice(a)
	require.Equal(t, e, a)
}

func TestLegacy(t *testing.T) {
	const (
		char = '&'
		l    = "  &eHallo, &oich &b&mbin &rstolz!"
	)

	require.Equal(t,
		`{"text":"  ","extra":[{"text":"Hallo, ","color":"yellow","extra":[{"text":"ich ","italic":true}]},{"text":"bin ","color":"aqua","strikethrough":true},{"text":"stolz!"}]}`,
		LegacyChar(l, char).String())
}

func TestLegacy2(t *testing.T) {
	var (
		c = DefaultLegacyChar
		l = "§cTest§l"
	)
	require.Equal(t,
		`{"text":"","extra":[{"text":"Test","color":"red","extra":[{"text":"","bold":true}]}]}`,
		LegacyChar(l, c).String())
}
