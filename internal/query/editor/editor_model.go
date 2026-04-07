package editor

import (
	"query-terminal/internal/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var submissionKey = key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "submit query"))

type UserInputMsg struct {
	InputQuery string
}

type Model struct {
	text textarea.Model
	help help.Model
}

func NewQueryEditorModel() Model {
	model := textarea.New()
	model.Focus()
	model.ShowLineNumbers = true
	model.Placeholder = "SELECT * FROM users;"
	model.SetWidth(50)
	model.SetHeight(20)

	return Model{
		text: model,
		help: help.New(),
	}
}

func (m Model) ShortHelp() []key.Binding {
	return getKeyBindings()
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{getKeyBindings()}
}

func getKeyBindings() []key.Binding {
	return []key.Binding{
		submissionKey,
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, submissionKey) {
			return m, func() tea.Msg {
				return UserInputMsg{
					InputQuery: m.text.Value(),
				}
			}
		}

	case tea.WindowSizeMsg:
		m.text.SetWidth(msg.Width - 10)
	}

	var cmd tea.Cmd
	m.text, cmd = m.text.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, m.text.View(), m.help.View(m))
}
