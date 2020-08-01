package legacy

import (
	. "go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/codec"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Legacy struct {
	Char    rune // The format char to use (e.g. DefaultChar, SectionChar, AmpersandChar).
	HexChar rune // The hex prefixing char to use (e.g. HexChar).

	// Since Minecraft 1.16+ there can be hex colors (e.g. "#ff5555" instead of the named color "red").
	// This setting decides whether to use the nearest legacy named color of the hex color instead.
	//
	// This setting is false by default to support older client versions.
	NoDownsampleColor bool

	// Whether to add a "open_url" click event with the URL to the text containing an URL.
	ClickableUrl bool
}

var _ codec.Codec = (*Legacy)(nil)

var urlRegex = regexp.MustCompile(`(?:(https?)://)?([-\w_.]+\.\w{2,})(/\S*)?`)

const (
	DefaultChar    = SectionChar
	DefaultHexChar = HexChar

	// The legacy character used by Minecraft.
	SectionChar rune = '§'
	// The legacy character frequently used by configurations and commands.
	AmpersandChar rune = '&'
	// The character used to prefix hex colors.
	HexChar rune = '#'

	Chars = "0123456789abcdefklmnor"
)

func (l *Legacy) Marshal(wr io.Writer, c Component) error {
	if l.Char == 0 {
		l.Char = DefaultChar
	}
	if l.HexChar == 0 {
		l.HexChar = DefaultHexChar
	}
	s := newStringBuilder(l, l.Char)
	s.append(c, &style{b: s, decorations: map[Decoration]struct{}{}})
	_, err := wr.Write([]byte(s.String()))
	return err
}

// Reset is the legacy reset format.
type Reset struct{}

// String implements component.Format.
func (Reset) String() string {
	return "reset"
}

var formats = func() (f []Format) {
	for _, n := range NamesOrder {
		f = append(f, n)
	}
	for _, d := range DecorationsOrder {
		f = append(f, d)
	}
	f = append(f, Reset{})

	if len(f) != len(Chars) { // assert same length
		panic("formats length differs from legacy.Chars length")
	}
	return
}()

func formatIndex(f Format) int {
	for i, test := range formats {
		if test == f {
			return i
		}
	}
	return -1
}

// legacy string builder
type stringBuilder struct {
	strings.Builder
	style *style
	char  rune

	l *Legacy
}

func newStringBuilder(l *Legacy, char rune) *stringBuilder {
	b := &stringBuilder{l: l, char: char}
	b.style = newStyle(b)
	return b
}

func (b *stringBuilder) append(c Component, s *style) {
	s.apply(c)

	if t, ok := c.(*Text); ok && len(t.Content) != 0 {
		s.applyFormat()
		_, _ = b.WriteString(t.Content)
	}

	if len(c.Children()) == 0 {
		return
	}
	childrenStyle := s.copy()
	for _, child := range c.Children() {
		b.append(child, childrenStyle)
		childrenStyle.set(s)
	}
}

func (b *stringBuilder) appendFormat(format Format) {
	_, _ = b.WriteRune(b.char)
	_, _ = b.WriteString(b.toLegacyCode(format))
}

// Code returns the legacy format.
func (b *stringBuilder) toLegacyCode(format Format) string {
	if color, ok := format.(Color); ok {
		if b.l.NoDownsampleColor {
			return color.Hex()
		} else {
			format = color.Named()
		}
	}
	return string(Chars[formatIndex(format)])
}

// legacy style format
type style struct {
	b *stringBuilder

	color       Color
	decorations decorations
}

func newStyle(b *stringBuilder) *style {
	return &style{b: b, decorations: map[Decoration]struct{}{}}
}

// This is a set.
type decorations map[Decoration]struct{}

// clear instead of creating new since
// this map is referenced by other code pieces
func (d decorations) clear() {
	for dec := range d {
		delete(d, dec)
	}
}

func (d decorations) hasAll(d2 decorations) bool {
	if len(d) != len(d2) {
		return false
	}
	for deco := range d2 {
		if !d.has(deco) {
			return false
		}
	}
	return true
}

func (d decorations) has(deco Decoration) bool {
	_, ok := d[deco]
	return ok
}

func (d decorations) addAll(d2 decorations) {
	for deco := range d2 {
		d[deco] = struct{}{}
	}
}

func (s *style) copy() *style {
	return &style{
		b:     s.b,
		color: s.color,
		decorations: func() decorations {
			m := make(decorations, len(s.decorations))
			for d := range s.decorations {
				m[d] = struct{}{}
			}
			return m
		}(),
	}
}

func (s *style) apply(c Component) {
	color := c.Style().Color
	if color != nil {
		s.color = color
	}

	for d := range Decorations {
		switch c.Style().Decoration(d) {
		case True:
			s.decorations[d] = struct{}{}
		case False:
			delete(s.decorations, d)
		}
	}
}

func (s *style) applyFormat() {
	// If color changes, we need to do a full reset
	if s.color != s.b.style.color {
		s.applyFullFormat()
		return
	}

	// Does current have any decorations we don't have?
	// Since there is no way to undo decorations, we need to reset these cases
	if !s.decorations.hasAll(s.b.style.decorations) {
		s.applyFullFormat()
		return
	}

	// Apply new decorations
	for d := range s.decorations {
		if s.b.style.decorations.has(d) {
			continue
		}
		s.b.style.decorations[d] = struct{}{}
		s.b.appendFormat(d)
	}
	return
}

func (s *style) set(s2 *style) {
	s.color = s2.color
	s.decorations.clear()
	s.decorations.addAll(s2.decorations)
}

func (s *style) applyFullFormat() {
	if s.color != nil {
		s.b.appendFormat(s.color)
	} else {
		s.b.appendFormat(Reset{})
	}
	s.b.style.color = s.color
	for d := range s.decorations {
		s.b.appendFormat(d)
	}

	s.b.style.decorations.clear()
	s.b.style.decorations.addAll(s.decorations)

}

// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode

// Unmarshal takes a string and always returns the *component.Text from it or an error.
func (l *Legacy) Unmarshal(data []byte) (Component, error) {
	if l.Char == 0 {
		l.Char = DefaultChar
	}
	if l.HexChar == 0 {
		l.HexChar = DefaultHexChar
	}
	input := string(data)

	next := lastIndexFrom(input, l.Char, len(input)-1)
	if next == -1 {
		return l.extractUrl(&Text{Content: input}), nil
	}

	var (
		parts   []Component
		current *Text
		reset   bool
		pos     = len(input)
	)
	for {
		_, format, ok := decodeFormat(rune(input[next+1]))
		if ok {
			from := next + 2
			if from != pos {
				if current == nil {
					current = &Text{}
				} else if reset {
					parts = append(parts, current)
					reset = false
					current = &Text{}
				} else {
					current = &Text{Extra: []Component{current}}
				}

				current.Content = onlyValidUTF8(input[from:pos])
			} else if current == nil {
				current = &Text{}
			}

			if !reset {
				reset = applyFormat(current, format)
			}
			pos = next
		}
		next = lastIndexFrom(input, l.Char, next)

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
	return &Text{Content: s, Extra: parts}, nil
}

func (l *Legacy) extractUrl(t *Text) Component {
	if !l.ClickableUrl {
		return t
	}
	url := urlRegex.FindString(t.Content)
	if len(url) != 0 {
		t.S.ClickEvent = OpenUrl(url)
	}
	return t
}

func applyFormat(t *Text, format Format) bool {
	switch f := format.(type) {
	case *RGB:
		t.S.Color = f
		return true
	case *Named:
		t.S.Color = f.RGB
		return true
	case Decoration:
		t.S.SetDecoration(f, True)
		return false
	case *Decoration:
		t.S.SetDecoration(*f, True)
		return false
	case Reset, *Reset:
		return true
	default:
		return false // unknown format, just returns false!
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

func onlyValidUTF8(s string) string {
	b := make([]rune, 0, len(s))
	for _, r := range s {
		if r != utf8.RuneError {
			b = append(b, r)
		}
	}
	return string(b)
}

func reverseSlice(a []Component) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

// returned values only valid if returns true
func decodeFormat(legacy rune) (t FormatCodeType, f Format, ok bool) {
	t, ok = determineFormatType(legacy)
	if !ok {
		return 0, nil, false
	}
	switch t {
	case MojangLegacy:
		if i := strings.IndexRune(Chars, legacy); i != -1 {
			return MojangLegacy, formats[i], true
		}
	}
	return 0, nil, false
}

func determineFormatType(char rune) (FormatCodeType, bool) {
	if strings.IndexRune(Chars, char) != -1 {
		return MojangLegacy, true
	}
	return 0, false
}

type FormatCodeType uint8

const (
	MojangLegacy FormatCodeType = iota
	// ...maybe add Bungeecord's bad format in the future (╯°□°)╯︵ ┻━┻
)
