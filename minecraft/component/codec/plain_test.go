package codec

import (
	"github.com/stretchr/testify/require"
	"go.minekube.com/common/minecraft/component"
	"strings"
	"testing"
)

var p = &Plain{}

func TestPlain_Marshal(t *testing.T) {
	b := new(strings.Builder)
	err := p.Marshal(b, txt)
	require.NoError(t, err)
	require.Equal(t, b.String(), "Hello there!")
}

func TestPlain_Unmarshal(t *testing.T) {
	c, err := p.Unmarshal([]byte("Hello there!"))
	require.NoError(t, err)
	tx, ok := c.(*component.Text)
	require.True(t, ok)
	require.Equal(t, tx, &component.Text{Content: "Hello there!"})
}
