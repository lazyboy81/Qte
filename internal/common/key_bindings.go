package common

import (
	"github.com/charmbracelet/bubbles/key"
)

var QuitKey = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("q", "quit"),
)

var StartKey = key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("enter", "start"),
)
var ReturnKey = key.NewBinding(
	key.WithKeys("backspace"),
	key.WithHelp("backspace", "return"),
)

type KeyMap struct {
	Quit   key.Binding
	Start  key.Binding
	Return key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Start, k.Return, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.Return, k.Quit},
	}
}

var mappings *KeyMap

func newMapping() {
	mappings = &KeyMap{
		Start:  StartKey,
		Return: ReturnKey,
		Quit:   QuitKey,
	}
}

func GetKeyMappings() KeyMap {
	if mappings == nil {
		newMapping()
	}
	return *mappings
}
