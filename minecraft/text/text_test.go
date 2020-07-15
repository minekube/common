package text

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestText(t *testing.T) {
	txt := Text{
		Text: "Hallo ",
		ClickEvent: &ClickEvent{
			Action: OpenUrl,
			Value:  "https://google.com/",
		},
		HoverEvent: &HoverEvent{
			Action: HoverEventShowText,
			Value:  &Text{Text: "Test"},
		},
		Italic: True,
		Extra: []*Text{
			{
				Text: "du da!",
			},
		},
	}
	j := `{"text":"Hallo ","italic":true,"clickEvent":{"action":"open_url","value":"https://google.com/"},"hoverEvent":{"action":"show_text","value":{"text":"Test"}},"extra":[{"text":"du da!"}]}`
	require.Equal(t, j, txt.String())
	var jTxt Text
	require.NoError(t, json.Unmarshal([]byte(j), &jTxt))
	require.Equal(t, txt, jTxt)
}

func TestState(t *testing.T) {
	txt := Text{
		Text: "Hello ",
		Bold: True,
		Extra: []*Text{
			{
				Text: " you!",
			},
		},
	}
	j := `{"text":"Hello ","bold":true,"extra":[{"text":" you!"}]}`
	require.Equal(t, j, txt.String())
	var jTxt Text
	require.NoError(t, json.Unmarshal([]byte(j), &jTxt))
	require.Equal(t, txt, jTxt)

	txt = Text{
		Text: "Hello ",
		Bold: True,
		Extra: []*Text{
			{
				Text: " you!",
				Bold: False,
			},
		},
	}
	j = `{"text":"Hello ","bold":true,"extra":[{"text":" you!","bold":false}]}`
	require.Equal(t, j, txt.String())
	require.NoError(t, json.Unmarshal([]byte(j), &jTxt))
	require.Equal(t, txt, jTxt)
}

func TestTranslation(t *testing.T) {
	tr := &Translation{
		Translate: "test.key",
		With: []Component{
			&Text{Text: "Hello"},
		},
	}
	s := `{"translate":"test.key","with":[{"text":"Hello"}]}`
	require.Equal(t, s, tr.String())
}
