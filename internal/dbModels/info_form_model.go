package dbModels

import (
	"errors"
	"qit/internal/common"
	"qit/internal/db"

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
						huh.NewOption[string](db.DbSqlite, db.DbSqlite),
						huh.NewOption[string](db.DbOracle, db.DbOracle),
						huh.NewOption[string](db.DbPostgres, db.DbPostgres),
					).
					Value(&selectedProvider).
					Inline(true),
				huh.NewInput().Title("DB URL").Value(&databaseUrl).
					Validate(func(value string) error {
						if len(value) == 0 {
							return errors.New("database url cannot be empty")
						}

						return nil
					}),
				huh.NewInput().Title("Username").Placeholder("Optional").Value(&username),
				huh.NewInput().Title("Password").Placeholder("Optional").Value(&password).EchoMode(huh.EchoModePassword),
			),
		),
	}
}

func (m InfoFormModel) EnterParams(params common.AdditionalParam) {}

func (m InfoFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m InfoFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.form.Update(msg)
	m.form = updated.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, common.SwitchTo(common.ScreenConnection,
			common.AdditionalParam{"provider": selectedProvider, "url": databaseUrl, "username": username, "password": password},
		)
	}

	return m, cmd
}

func (m InfoFormModel) View() string {
	return m.form.View()
}
