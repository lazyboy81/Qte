package dbModels

import (
	"qte/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	isTransactional bool
	query           string
)

type QueryModel struct {
	form *huh.Form
}

func NewQueryModel() QueryModel {
	return QueryModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Is Update/Insert/Delete?").
					Affirmative("Yes").
					Negative("No").
					Value(&isTransactional),
				huh.NewText().
					Title("SQL Query").
					ShowLineNumbers(true).
					Placeholder("SELECT * FROM users;").
					Lines(8). // Set height
					Value(&query).
					Editor("vim").
					ExternalEditor(true),
			),
		),
	}
}

func (m QueryModel) EnterParams(params common.AdditionalParam) {}

func (m QueryModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m QueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.form.Update(msg)
	m.form = updated.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, common.SwitchTo(common.QueryResultScreen, nil)
	}

	return m, cmd
}

func (m QueryModel) View() string {
	return m.form.View()
}
