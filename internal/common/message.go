package common

import tea "github.com/charmbracelet/bubbletea"

// SwitchScreenMsg Navigation message — any model can emit this to switch screens
type SwitchScreenMsg int

const (
	ScreenHome SwitchScreenMsg = iota
	ScreenInfo
	ScreenConnection
)

// ErrMsg Generic error
type ErrMsg struct {
	Err error
}

// SwitchTo Helper to create a command that switches screens
func SwitchTo(screen SwitchScreenMsg) tea.Cmd {
	return func() tea.Msg {
		return screen
	}
}
