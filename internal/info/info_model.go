package info

import (
	"TQT/internal/common"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#874BFD"))
)

type Model struct{}

func NewInfoModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.GetBinding("backspace")):
			return m, func() tea.Msg {
				return common.ScreenHome
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	body := `
Welcome to the Terminal Query Tool application.

This tool helps you run queries from the comfort of your terminal.
`

	content := boxStyle.Render(body)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		content,
	)
}
