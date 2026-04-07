package info

import (
	"query-terminal/internal/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
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

func (m Model) EnterParams(params common.AdditionalParam) {}

func NewInfoModel() Model {
	return Model{
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlR.String():
			return m, common.SwitchTo(common.HomeScreen, nil)
		}
	}

	return m, nil
}

func (m Model) View() string {
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

	return lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		m.help.View(KeyMap{
			Return: key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "return")),
		}),
	)
}
