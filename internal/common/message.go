package common

import (
	"cmp"

	tea "charm.land/bubbletea/v2"
)

// SwitchScreenMsg Navigation message — any model can emit this to switch screens
type SwitchScreenMsg int

type ScreenMsg struct {
	DestScreen SwitchScreenMsg
}

type ConnectionInfo struct {
	Driver string
	Url    string
	Name   string
}

const (
	HomeScreen SwitchScreenMsg = iota
	InfoScreen
	ConnectionFormScreen
	ConnectionEstablishmentScreen
	ConnectionPickerScreen
	QueryScreen
)

// ErrMsg Generic error
type ErrMsg struct {
	Err error
}

// SwitchTo Helper to create a command that switches screens
func SwitchTo(screen SwitchScreenMsg) tea.Cmd {
	return func() tea.Msg {
		return ScreenMsg{
			DestScreen: screen,
		}
	}
}

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
