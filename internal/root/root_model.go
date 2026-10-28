package root

import (
	"query-terminal/internal/common"
	"query-terminal/internal/home"
	"query-terminal/internal/info"
	"strings"

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
	//screens       map[common.SwitchScreenMsg]tea.Model
	tabs          []tea.Model
	activeTab     int
	width, height int
	titleStyle    lipgloss.Style
}

func NewRootModel() *Model {
	return &Model{
		currentScreen: common.HomeScreen,
		//screens: map[common.SwitchScreenMsg]tea.Model{
		//	common.HomeScreen:                    home.GetHomePage(),
		//	common.InfoScreen:                    info.GetInfoModel(),
		//	common.ConnectionFormScreen:          form.GetConnectionFormModel(),
		//	common.ConnectionPickerScreen:        picker.GetConnectionPickerModel(),
		//	common.ConnectionEstablishmentScreen: establish.GetConnectionEstablishmentModel(),
		//	common.QueryScreen:                   query.GetQueryModel(),
		//},
		tabs: []tea.Model{
			home.GetHomePage(),
			info.GetInfoModel(),
			//form.GetConnectionFormModel(),
			//picker.GetConnectionPickerModel(),
			//establish.GetConnectionEstablishmentModel(),
			//query.GetQueryModel(),
		},
		titleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).PaddingBottom(1),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
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
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	// Handle window resize globally, store dimensions
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to current screen too (fall through below)

		// *** THIS IS THE KEY PART — screen switching ***
		//case common.ScreenMsg:
		//	m.currentScreen = msg.DestScreen
		//	// Re-initialize the new screen when switching to it
		//	return m, m.screens[m.currentScreen].Init()
		//}
	}

	// Delegate to the current screen
	//updatedScreen, cmd := m.screens[m.currentScreen].Update(msg)
	//m.screens[m.currentScreen] = updatedScreen

	return m, nil
}

func (m *Model) View() tea.View {

	doc := strings.Builder{}
	var renderedTabs []string

	for _, t := range m.tabs {
		renderedTabs = append(renderedTabs, t.View().Content)
	}

	//rendered := lipgloss.JoinVertical(
	//	lipgloss.Center,
	//	m.titleStyle.Render(logo),
	//	m.screens[m.currentScreen].View().Content,
	//)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	place := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, doc.String())

	view := tea.NewView(place)
	view.AltScreen = true
	view.WindowTitle = "Qte"
	return view
}
