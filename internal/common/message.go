package common

import tea "github.com/charmbracelet/bubbletea"

// SwitchScreenMsg Navigation message — any model can emit this to switch screens
type SwitchScreenMsg int

type AdditionalParam map[string]any

type ScreenMsg struct {
	DestScreen SwitchScreenMsg
	Params     AdditionalParam
}

const (
	HomeScreen SwitchScreenMsg = iota
	InfoScreen
	ConnectionFormScreen
	ConnectionEstablishmentScreen
	ConnectionPickerScreen
	QueryEditorScreen
	QueryResultScreen
)

// ErrMsg Generic error
type ErrMsg struct {
	Err error
}

// SwitchTo Helper to create a command that switches screens
func SwitchTo(screen SwitchScreenMsg, params AdditionalParam) tea.Cmd {
	return func() tea.Msg {
		return ScreenMsg{
			DestScreen: screen,
			Params:     params,
		}
	}
}
