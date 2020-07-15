package text

import (
	"strings"
	"unicode/utf8"
)

const (
	DefaultLegacyChar = '§'
	LegacyChars       = "0123456789abcdefklmnor"
)

var (
	decorations = []string{
		string(Obfuscated),
		string(Bold),
		string(Strikethrough),
		string(Underlined),
		string(Italic),
	}
	colors = []string{
		string(Black),
		string(DarkBlue),
		string(DarkGreen),
		string(DarkAqua),
		string(DarkRed),
		string(DarkPurple),
		string(Gold),
		string(Gray),
		string(DarkGray),
		string(Blue),
		string(Green),
		string(Aqua),
		string(Red),
		string(LightPurple),
		string(Yellow),
		string(White),
	}
	formats []string

	decorationSet = newSet(decorations...)
)

func init() {
	formats = make([]string, 0, len(LegacyChars))
	formats = append(formats, colors...)
	formats = append(formats, decorations...)
	formats = append(formats, Reset)
}

func formatByLegacyChar(legacy rune) string {
	i := strings.IndexByte(LegacyChars, byte(legacy))
	return formats[i]
}

// Legacy converts a legacy formatted chat string with the
// § character to Text.
func Legacy(input string) *Text {
	return LegacyChar(input, DefaultLegacyChar)
}

// Legacy converts a legacy formatted chat string with a given character to Text.
func LegacyChar(input string, char rune) *Text {
	next := lastIndexFrom(input, char, len(input)-1)
	if next == -1 {
		return &Text{Text: input}
	}

	var (
		parts   []*Text
		current *Text
		reset   bool
		pos     = len(input)
	)
	for {
		format := formatByLegacyChar(rune(input[next+1]))
		if len(format) != 0 {
			from := next + 2
			if from != pos {
				if current == nil {
					current = &Text{}
				} else if reset {
					parts = append(parts, current)
					reset = false
					current = &Text{}
				} else {
					current = &Text{Extra: []*Text{current}}
				}

				current.Text = onlyValidUTF8(input[from:pos])
			} else if current == nil {
				current = &Text{}
			}

			reset = reset || applyFormat(current, format)
			pos = next
		}
		next = lastIndexFrom(input, char, next)

		if next == -1 {
			break
		}
	}

	if current != nil {
		parts = append(parts, current)
	}
	reverseSlice(parts)

	var s string
	if pos > 0 {
		s = onlyValidUTF8(input[:pos])
	}
	return &Text{Text: s, Extra: parts}
}

func onlyValidUTF8(s string) string {
	b := make([]rune, 0, len(s))
	for _, r := range s {
		if r != utf8.RuneError {
			b = append(b, r)
		}
	}
	return string(b)
}

func applyFormat(t *Text, format string) bool {
	if format == Reset {
		t.Reset()
		return true
	}
	if decorationSet.has(format) {
		switch Decoration(format) {
		case Italic:
			t.Italic = True
		case Underlined:
			t.Underlined = True
		case Strikethrough:
			t.Strikethrough = True
		case Bold:
			t.Bold = True
		case Obfuscated:
			t.Obfuscated = True
		}
		return false
	}
	t.Color = Color(format)
	return true
}

func reverseSlice(a []*Text) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func lastIndexFrom(s string, char rune, from int) int {
	from = min(len(s), from)
	if from < 0 {
		return -1
	}
	return strings.LastIndexByte(s[:from], byte(char))
}

func min(x, y int) int {
	if y < x {
		return y
	}
	return x
}

type stringSet map[string]struct{}

func newSet(s ...string) stringSet {
	set := stringSet{}
	for _, v := range s {
		set[v] = struct{}{}
	}
	return set
}

func (set stringSet) has(s string) bool {
	_, ok := set[s]
	return ok
}
