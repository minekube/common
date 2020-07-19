package component

var (
	// A decoration which makes text obfuscated/unreadable.
	Obfuscated TextDecoration = "obfuscated"
	// A decoration which makes text appear bold.
	Bold TextDecoration = "bold"
	// A decoration which makes text have a strike through it.
	Strikethrough TextDecoration = "strikethrough"
	// A decoration which makes text have an underline.
	Underlined TextDecoration = "underlined"
	// A decoration which makes text appear in italics.
	Italic TextDecoration = "italic"

	Decorations = func() map[TextDecoration]struct{} {
		m := map[TextDecoration]struct{}{}
		for _, a := range []TextDecoration{
			Obfuscated,
			Bold,
			Strikethrough,
			Underlined,
			Italic,
		} {
			m[a] = struct{}{}
		}
		return m
	}()
)

type TextDecoration string
