package info

import (
	"query-terminal/internal/common"
	"sync"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/list"
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			BorderForeground(lipgloss.Color("#874BFD"))

	listStyle = lipgloss.NewStyle().Padding(1, 2)
)

type KeyMap struct {
	Return key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Return}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Return},
	}
}

type Model struct {
	help help.Model
}

var model *Model
var once sync.Once

func GetInfoModel() *Model {
	once.Do(func() {
		model = newInfoModel()
	})

	return model
}

func newInfoModel() *Model {
	return &Model{
		help: help.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return m, common.SwitchTo(common.HomeScreen)
		}
	}

	return m, nil
}

func (m *Model) View() tea.View {
	body := `
Welcome to InputQuery; in Terminal application.
This tool helps you run queries from the comfort of your terminal.

Supported database include: 
`

	supportedDbs := list.New().
		Item("Oracle").
		Item("SQLite").
		Item("PostgreSQL").
		Enumerator(list.Bullet)

	content := boxStyle.Render(body + listStyle.Render(supportedDbs.String()))

	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			m.help.View(KeyMap{
				Return: key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "return")),
			}),
		),
	)
}
