package root

import (
	"query-terminal/internal/common"
	"query-terminal/internal/connection/establish"
	"query-terminal/internal/connection/form"
	"query-terminal/internal/connection/picker"
	"query-terminal/internal/home"
	"query-terminal/internal/info"
	"query-terminal/internal/query"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	logo = `
   ____   ______        ;-)
  / __ \ /_  __/___     
 / / / /  / /  / _ \    
/ /_/ /  / /  /  __/    
\___\_\ /_/   \___/
`
)

type Model struct {
	currentScreen common.SwitchScreenMsg
	screens       map[common.SwitchScreenMsg]tea.Model
	width, height int
	titleStyle    lipgloss.Style
}

func NewRootModel() *Model {
	return &Model{
		currentScreen: common.HomeScreen,
		screens: map[common.SwitchScreenMsg]tea.Model{
			common.HomeScreen:                    home.GetHomePage(),
			common.InfoScreen:                    info.GetInfoModel(),
			common.ConnectionFormScreen:          form.GetConnectionFormModel(),
			common.ConnectionPickerScreen:        picker.GetConnectionPickerModel(),
			common.ConnectionEstablishmentScreen: establish.GetConnectionEstablishmentModel(),
			common.QueryScreen:                   query.GetQueryModel(),
		},
		titleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).PaddingBottom(1),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.screens[m.currentScreen].Init()
}

// Update we switch on the msg type that was given to us and match it
// against the key mappings that we have defined
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, common.SwitchTo(common.InfoScreen)
		}
	// Handle window resize globally, store dimensions
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	// Forward to current screen too (fall through below)

	// *** THIS IS THE KEY PART — screen switching ***
	case common.ScreenMsg:
		m.currentScreen = msg.DestScreen
		// Re-initialize the new screen when switching to it
		return m, m.screens[m.currentScreen].Init()
	}

	// Delegate to the current screen
	updatedScreen, cmd := m.screens[m.currentScreen].Update(msg)
	m.screens[m.currentScreen] = updatedScreen

	return m, cmd
}

func (m *Model) View() tea.View {
	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Center,
		m.titleStyle.Render(logo),
		m.screens[m.currentScreen].View().Content,
	))
}
