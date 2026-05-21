package codec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.minekube.com/common/minecraft/component"
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

func TestPlain_Marshal_score(t *testing.T) {
	b := new(strings.Builder)
	err := p.Marshal(b, &component.Score{
		Name:      "PlayerOne",
		Objective: "kills",
		Value:     "42",
		Extra: []component.Component{
			&component.Text{Content: " points"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, " points", b.String())
}

func TestPlain_Marshal_scoreSkipped_AppendedTextVisible(t *testing.T) {
	b := new(strings.Builder)
	err := p.Marshal(b, &component.Text{Extra: []component.Component{
		&component.Score{Name: "hard", Objective: "true"},
		&component.Text{Content: "Hello world"},
	}})
	require.NoError(t, err)
	require.Equal(t, "Hello world", b.String())
}
