package component

type ClickEvent interface {
	Action() ClickAction
	Value() string
}

func NewClickEvent(action ClickAction, value string) ClickEvent {
	return &clickEvent{action, value}
}

type clickEvent struct {
	action ClickAction
	value  string
}

func (c *clickEvent) Action() ClickAction {
	return c.action
}
func (c *clickEvent) Value() string {
	return c.value
}

type ClickAction interface {
	Name() string   // The ClickAction name.
	Readable() bool // When an clickAction is not readable it will not be unmarshalled.
}

type clickAction struct {
	name     string
	readable bool
}

func (a *clickAction) Name() string {
	return a.name
}

func (a *clickAction) Readable() bool {
	return a.readable
}

func (a *clickAction) String() string {
	return a.name
}

func OpenUrl(url string) ClickEvent {
	return &clickEvent{OpenUrlAction, url}
}

func OpenFile(file string) ClickEvent {
	return &clickEvent{OpenFileAction, file}
}

func RunCommand(command string) ClickEvent {
	return &clickEvent{RunCommandAction, command}
}

func SuggestCommand(command string) ClickEvent {
	return &clickEvent{SuggestCommandAction, command}
}

func ChangePage(page string) ClickEvent {
	return &clickEvent{ChangePageAction, page}
}

func CopyToClipboard(text string) ClickEvent {
	return &clickEvent{CopyToClipboardAction, text}
}

func ShowDialog(dialog string) ClickEvent {
	return &clickEvent{ShowDialogAction, dialog}
}

func CustomEvent(id string, payload ...string) ClickEvent {
	value := id
	if len(payload) > 0 && payload[0] != "" {
		value = id + "|" + payload[0]
	}
	return &clickEvent{CustomEventAction, value}
}

var (
	OpenUrlAction         ClickAction = &clickAction{"open_url", true}
	OpenFileAction        ClickAction = &clickAction{"open_file", false}
	RunCommandAction      ClickAction = &clickAction{"run_command", true}
	SuggestCommandAction  ClickAction = &clickAction{"suggest_command", true}
	ChangePageAction      ClickAction = &clickAction{"change_page", true}
	CopyToClipboardAction ClickAction = &clickAction{"copy_to_clipboard", true}
	ShowDialogAction      ClickAction = &clickAction{"show_dialog", true}
	CustomEventAction     ClickAction = &clickAction{"custom", true}

	ClickActions = func() map[string]ClickAction {
		m := map[string]ClickAction{}
		for _, a := range []ClickAction{
			OpenUrlAction,
			OpenFileAction,
			RunCommandAction,
			SuggestCommandAction,
			ChangePageAction,
			CopyToClipboardAction,
			ShowDialogAction,
			CustomEventAction,
		} {
			m[a.Name()] = a
		}
		return m
	}()
)
