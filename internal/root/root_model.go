package root

import (
	"query-terminal/internal/common"
	"query-terminal/internal/connection/establish"
	"query-terminal/internal/connection/form"
	"query-terminal/internal/connection/picker"
	"query-terminal/internal/home"
	"query-terminal/internal/info"
	"query-terminal/internal/query"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// because each one returns a new Style struct we cannot use const keyword
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).PaddingBottom(1)

type Screen interface {
	tea.Model

	EnterParams(common.AdditionalParam)
}

type Model struct {
	currentScreen common.SwitchScreenMsg
	screens       map[common.SwitchScreenMsg]Screen
	width, height int
}

func NewRootModel() Model {
	return Model{
		currentScreen: common.HomeScreen,
		screens: map[common.SwitchScreenMsg]Screen{
			common.HomeScreen:                    home.NewHomePage(),
			common.InfoScreen:                    info.NewInfoModel(),
			common.ConnectionFormScreen:          form.NewConnectionFormModel(),
			common.ConnectionPickerScreen:        picker.NewConnectionPickerModel(),
			common.ConnectionEstablishmentScreen: establish.NewConnectionEstablishmentModel(),
			common.QueryScreen:                   query.NewQueryModel(),
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
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case "?":
			return m, common.SwitchTo(common.InfoScreen, nil)
		}
	// Handle window resize globally, store dimensions
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	// Forward to current screen too (fall through below)

	// *** THIS IS THE KEY PART — screen switching ***
	case common.ScreenMsg:
		m.currentScreen = msg.DestScreen
		m.screens[m.currentScreen].EnterParams(msg.Params)
		// Re-initialize the new screen when switching to it
		return m, m.screens[m.currentScreen].Init()
	}

	// Delegate to the current screen
	updatedScreen, cmd := m.screens[m.currentScreen].Update(msg)
	m.screens[m.currentScreen] = updatedScreen.(Screen)

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(logo),
		m.screens[m.currentScreen].View(),
	)
}
