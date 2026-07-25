package query

import (
	"fmt"
	"query-terminal/internal/db"
	"query-terminal/internal/query/editor"
	"query-terminal/internal/query/result"
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Tab struct {
	editor *editor.Model
	result *result.Model
}

func (t *Tab) Init() tea.Cmd {
	return t.result.Init() // since editor does not need any initialization
}

type Model struct {
	currentTab     int
	tabs           []*Tab
	titleStyle     lipgloss.Style
	boxStyle       lipgloss.Style
	tabStyle       lipgloss.Style
	activeTabStyle lipgloss.Style
}

var model *Model
var once sync.Once

func GetQueryModel() *Model {
	once.Do(func() {
		model = newQueryModel()
	})

	return model
}

func newQueryModel() *Model {
	tabs := make([]*Tab, 1)
	tabs[0] = &Tab{
		editor: editor.NewQueryEditorModel(),
		result: result.NewResultModel(),
	}

	tabStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(lipgloss.Color("#777"))
	activeTabStyle := tabStyle.Bold(true).Foreground(lipgloss.Color("#00ffcc")).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#00ffcc"))

	return &Model{
		currentTab:     0,
		tabs:           tabs,
		titleStyle:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")),
		boxStyle:       lipgloss.NewStyle().Width(60).Height(20).Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("#874BFD")),
		tabStyle:       tabStyle,
		activeTabStyle: activeTabStyle,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.tabs[0].Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+l" {
			m.cycleTabsForward()
		}

		if msg.String() == "ctrl+h" {
			m.cycleTabsBackward()
		}

		if msg.String() == "ctrl+t" {
			cmds = append(cmds, db.ConnectToDatabase("", "", "", len(m.tabs)+1))
		}

	case *db.ActiveConnection:
		m.tabs[msg.ConnId].result.SetConnection(msg)

	case tea.WindowSizeMsg:
		m.boxStyle.Width(msg.Width - m.boxStyle.GetHorizontalBorderSize())
	}

	updatedEditor, cmd := m.tabs[m.currentTab].editor.Update(msg)
	cmds = append(cmds, cmd)
	m.tabs[m.currentTab].editor = updatedEditor.(*editor.Model)

	updatedResult, cmd := m.tabs[m.currentTab].result.Update(msg)
	cmds = append(cmds, cmd)
	m.tabs[m.currentTab].result = updatedResult.(*result.Model)

	return m, tea.Sequence(cmds...)
}

func (m *Model) View() tea.View {
	var tabs string
	for i := range m.tabs {

		tabName := fmt.Sprintf("query window %d", i)

		if i == m.currentTab {
			tabs += m.activeTabStyle.Render(tabName)
		} else {
			tabs += m.tabStyle.Render(tabName)
		}
	}

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Center,
			tabs,
			lipgloss.JoinVertical(
				lipgloss.Center,
				m.titleStyle.Align(lipgloss.Left).Render("Editor"),
				m.boxStyle.Render(m.tabs[m.currentTab].editor.View().Content),
				m.titleStyle.Align(lipgloss.Left).Render("Query Results"),
				m.boxStyle.Render(m.tabs[m.currentTab].result.View().Content),
			),
			"(ctrl+h/ctrl+l to switch, ctrl+t to create new tab,  ctrl+c to quit)",
		),
	)

}

func (m *Model) cycleTabsForward() {
	m.currentTab = (m.currentTab + 1) % len(m.tabs)
}

func (m *Model) cycleTabsBackward() {
	temp := m.currentTab - 1
	if temp < 0 {
		temp = len(m.tabs) - 1
	}

	m.currentTab = temp
}
