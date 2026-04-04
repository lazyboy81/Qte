package editor

import (
	"query-terminal/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	query = "query"
)

type Model struct {
	form *huh.Form
}

func NewQueryEditorModel() Model {
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("SQL Query").
					ShowLineNumbers(true).
					Placeholder("SELECT * FROM users;").
					Lines(8). // Set height
					Key(query).
					Editor("vim").
					ExternalEditor(true),
			),
		),
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.form.Update(msg)
	m.form = updated.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, common.SwitchTo(common.QueryResultScreen, common.AdditionalParam{"query": m.form.GetString(query)})
	}

	return m, cmd
}

func (m Model) View() string {
	return m.form.View()
}
