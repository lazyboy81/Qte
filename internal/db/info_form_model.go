package db

import (
	"qit/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	selectedProvider string
	databaseUrl      string
	username         string
	password         string
)

type InfoFormModel struct {
	form *huh.Form
}

func NewConnectionInfoFormModel() InfoFormModel {
	return InfoFormModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Database of choice?").
					Options(
						huh.NewOption[string](dbSqlite, dbSqlite),
						huh.NewOption[string](dbOracle, dbOracle),
						huh.NewOption[string](dbPostgres, dbPostgres),
					).
					Value(&selectedProvider).
					Inline(true),
				huh.NewInput().Title("DB URL").Value(&databaseUrl),
				huh.NewInput().Title("Username").Value(&username),
				huh.NewInput().Title("Password").Value(&password),
			),
		),
	}
}

func (m InfoFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m InfoFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.form.Update(msg)
	m.form = updated.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, common.SwitchTo(common.ScreenConnection)
	}

	return m, cmd
}

func (m InfoFormModel) View() string {
	return m.form.View()
}
