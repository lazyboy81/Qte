package root

import (
	"TQT/internal/common"
	"TQT/internal/connection"
	"TQT/internal/home"
	"TQT/internal/info"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	logo = `
  _______   ____   _______
 |__   __| / __ \ |__   __|
    | |   | |  | |   | |
    | |   | |  | |   | |
    | |   | |__| |   | |
    |_|    \___\_\   |_|
`
)

type Model struct {
	currentScreen common.SwitchScreenMsg
	screens       map[common.SwitchScreenMsg]tea.Model
	help          help.Model
	keys          common.KeyMap
	width, height int
}

func NewRootModel() Model {
	return Model{
		currentScreen: common.ScreenHome,
		help:          help.New(),
		keys:          common.GetKeyMappings(),
		screens: map[common.SwitchScreenMsg]tea.Model{
			common.ScreenHome:       home.NewHomePage(),
			common.ScreenInfo:       info.NewInfoModel(),
			common.ScreenConnection: connection.NewConnectionModel(),
		},
	}
}

func (m Model) Init() tea.Cmd {
	return m.screens[m.currentScreen].Init()
}

// Update we switch on the msg type that was given to us and match it
// against the key mappings that we have defined
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			return m, func() tea.Msg {
				return common.ScreenInfo
			}
		}
	// Handle window resize globally, store dimensions
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	// Forward to current screen too (fall through below)

	// *** THIS IS THE KEY PART — screen switching ***
	case common.SwitchScreenMsg:
		m.currentScreen = msg
		// Re-initialize the new screen when switching to it
		return m, m.screens[m.currentScreen].Init()

	}

	// Delegate to the current screen
	updatedScreen, cmd := m.screens[m.currentScreen].Update(msg)
	m.screens[m.currentScreen] = updatedScreen.(tea.Model)

	return m, cmd
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		PaddingBottom(1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(logo),
		m.screens[m.currentScreen].View(),
		m.help.View(m.keys),
	)
}
