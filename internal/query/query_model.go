package query

import (
	"query-terminal/internal/common"
	"query-terminal/internal/query/editor"
	"query-terminal/internal/query/result"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	editor     editor.Model
	result     result.Model
	titleStyle lipgloss.Style
	boxStyle   lipgloss.Style
}

func NewQueryModel() Model {
	return Model{
		editor:     editor.NewQueryEditorModel(),
		result:     result.NewResultModel(),
		titleStyle: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")),
		boxStyle:   lipgloss.NewStyle().Width(60).Height(20).Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("#874BFD")),
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(m.editor.Init(), m.result.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.boxStyle.Width(msg.Width - m.boxStyle.GetHorizontalBorderSize())
	}

	var cmds []tea.Cmd

	updatedEditor, cmd := m.editor.Update(msg)
	cmds = append(cmds, cmd)
	m.editor = updatedEditor.(editor.Model)

	updatedResult, cmd := m.result.Update(msg)
	cmds = append(cmds, cmd)
	m.result = updatedResult.(result.Model)

	return m, tea.Sequence(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.titleStyle.Align(lipgloss.Left).Render("Editor"),
		m.boxStyle.Render(m.editor.View()),
		m.titleStyle.Align(lipgloss.Left).Render("Query Results"),
		m.boxStyle.Render(m.result.View()),
	)
}
