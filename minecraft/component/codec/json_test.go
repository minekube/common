package codec

import (
	"github.com/stretchr/testify/require"
	. "go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"io/ioutil"
	"strings"
	"testing"
)

var j = &Json{
	NoLegacyHover:     true,
	NoDownsampleColor: true,
	StdJson:           true,
}

func BenchmarkJson_Marshal(b *testing.B) {
	insertion := "insert me"
	tx := &Text{Content: "Hello", S: Style{
		Obfuscated:    True,
		Bold:          False,
		Strikethrough: NotSet,
		Underlined:    True,
		Italic:        False,
		Font:          DefaultFont,
		Color:         Aqua,
		ClickEvent:    SuggestCommand("/help"),
		HoverEvent: ShowText(&Text{
			Content: " world",
			Extra: []Component{
				&Text{Content: "!"},
			},
		}),
		Insertion: &insertion,
	}}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := j.Marshal(ioutil.Discard, tx)
		require.NoError(b, err)
	}
}

var (
	insert = "insert me"
	txt    = &Text{
		Content: "Hello",
		Extra: []Component{
			&Text{Content: " there!", S: Style{
				Color:      Red.RGB,
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
			Color:         Aqua.RGB,
			ClickEvent:    SuggestCommand("/help"),
			HoverEvent: ShowText(&Text{
				Content: " world",
				Extra: []Component{
					&Text{Content: "!"},
				},
			}),
			Insertion: &insert,
		}}
	jsonTxt = `{"bold":false,"clickEvent":{"action":"suggest_command","value":"/help"},"color":"#55ffff","extra":[{"color":"#ff5555","italic":true,"obfuscated":false,"text":" there!"}],"font":"minecraft:default","hoverEvent":{"action":"show_text","contents":{"extra":[{"text":"!"}],"text":" world"}},"insertion":"insert me","italic":false,"obfuscated":true,"text":"Hello","underlined":true}`
)

func TestJson_Marshal_text(t *testing.T) {
	b := new(strings.Builder)
	err := j.Marshal(b, txt)
	require.NoError(t, err)
	require.Equal(t, jsonTxt, b.String())
}

func TestJson_Unmarshal_text(t *testing.T) {
	c, err := j.Unmarshal([]byte(jsonTxt))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

func TestJson_translation(t *testing.T) {
	tr := &Translation{
		Key: "sample.key",
		S:   Style{Color: Red.RGB},
		With: []Component{
			&Text{
				Content: "Hello",
			},
			&Translation{
				Key: "another.key",
				S:   Style{},
			},
		},
	}
	s := new(strings.Builder)
	require.NoError(t, j.Marshal(s, tr))
	const exp = `{"color":"#ff5555","translate":"sample.key","with":[{"text":"Hello"},{"translate":"another.key"}]}`
	require.Equal(t, exp, s.String())

	tr2, err := j.Unmarshal([]byte(exp))
	require.NoError(t, err)
	require.Equal(t, tr, tr2)
}
