package form

import (
	"errors"
	"query-terminal/internal/common"
	"query-terminal/internal/db"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

const (
	selectedDriver = "driver"
	databaseUrl    = "url"
	dbName         = "name"
)

type Model struct {
	form *huh.Form
}

func NewConnectionFormModel() Model {
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Database of choice? *").
					Options(
						huh.NewOption[string](db.DbSqlite, "sqlite3"),
						huh.NewOption[string](db.DbOracle, db.DbOracle),
						huh.NewOption[string](db.DbPostgres, db.DbPostgres),
					).
					Key(selectedDriver).
					Inline(true),
				huh.NewInput().Title("DB URL *").Key(databaseUrl).
					Validate(func(value string) error {
						if len(value) == 0 {
							return errors.New("database url cannot be empty")
						}

						return nil
					}),
				huh.NewInput().Title("Name *").Key(dbName).
					Validate(func(value string) error {
						if len(value) == 0 {
							return errors.New("database name cannot be empty")
						}

						return nil
					}),
				huh.NewNote().Description("\\* fields are mandatory"),
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
		return m, common.SwitchTo(common.ConnectionEstablishmentScreen,
			common.AdditionalParam{"driver": m.form.GetString(selectedDriver), "url": m.form.GetString(databaseUrl), "name": m.form.GetString(dbName)},
		)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.form.View()
}
