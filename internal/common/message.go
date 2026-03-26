package common

import tea "github.com/charmbracelet/bubbletea"

// SwitchScreenMsg Navigation message — any model can emit this to switch screens
type SwitchScreenMsg int

type ScreenMsg struct {
	DestScreen SwitchScreenMsg
	Params     map[string]string
}

const (
	ScreenHome SwitchScreenMsg = iota
	ScreenInfo
	ScreenConnectionForm
	ScreenConnection
	ScreenQuery
)

// ErrMsg Generic error
type ErrMsg struct {
	Err error
}

// SwitchTo Helper to create a command that switches screens
func SwitchTo(screen SwitchScreenMsg, params map[string]string) tea.Cmd {
	return func() tea.Msg {
		return ScreenMsg{
			DestScreen: screen,
			Params:     params,
		}
	}
}
