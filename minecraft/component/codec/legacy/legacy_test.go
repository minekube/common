package legacy

import (
	"github.com/stretchr/testify/require"
	. "go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"strings"
	"testing"
)

var l = &Legacy{
	Char:              DefaultChar,
	HexChar:           DefaultHexChar,
	NoDownsampleColor: false,
	ClickableUrl:      false,
	UrlStyle:          nil,
}
var (
	insert = "insert me"
	txt    = &Text{
		Content: "Hello",
		Extra: []Component{
			&Text{Content: " there!", S: Style{
				Color:      &RedColor,
				Italic:     True,
				Obfuscated: False,
			}},
		},
		S: Style{
			Obfuscated:    True,
			Bold:          False,
			Strikethrough: NotSet,
			Underlined:    True,
			Italic:        False,
			Font:          DefaultFont,
			Color:         &AquaColor,
			ClickEvent:    SuggestCommand("/help"),
			HoverEvent: ShowText(&Text{
				Content: " world",
				Extra: []Component{
					&Text{Content: "!"},
				},
			}),
			Insertion: &insert,
		}}
)

func TestLegacy_Marshal(t *testing.T) {
	b := new(strings.Builder)
	err := l.Marshal(b, txt)
	require.NoError(t, err)

	// some format order might range since ranging through decoration map drops random elements
	require.Contains(t, []string{
		"§b§n§kHello§c§n§o there!",
		"§b§n§kHello§c§o§n there!",
		"§b§k§nHello§c§n§o there!",
		"§b§k§nHello§c§o§n there!",
	}, b.String(), "%q invalid", b.String())
}

func TestLegacy_Unmarshal(t *testing.T) {
	c, err := l.Unmarshal([]byte("§b§k§nHello§c§n§o there!"))
	require.NoError(t, err)

	b := new(strings.Builder)
	err = l.Marshal(b, c)
	require.NoError(t, err)

	// some format order might range since ranging through decoration map drops random elements
	require.Contains(t, []string{
		"§b§n§kHello§c§n§o there!",
		"§b§n§kHello§c§o§n there!",
		"§b§k§nHello§c§n§o there!",
		"§b§k§nHello§c§o§n there!",
	}, b.String(), "%q invalid", b.String())
}
