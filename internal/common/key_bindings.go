package common

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
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
		Start: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "start"),
		),
		Return: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "return"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func GetBinding(binding string) key.Binding {
	if slices.Contains(mappings.Start.Keys(), binding) {
		return mappings.Start
	} else if slices.Contains(mappings.Return.Keys(), binding) {
		return mappings.Return
	} else {
		return mappings.Quit
	}
}

func GetKeyMappings() KeyMap {
	if mappings == nil {
		newMapping()
	}
	return *mappings
}
