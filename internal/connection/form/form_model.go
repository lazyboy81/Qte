package form

import (
	"errors"
	"query-terminal/internal/common"
	"query-terminal/internal/db"
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

const (
	selectedDriver = "driver"
	databaseUrl    = "url"
	dbName         = "name"
)

type Model struct {
	form *huh.Form
}

var model *Model
var once sync.Once

func GetConnectionFormModel() *Model {
	once.Do(func() {
		model = newConnectionFormModel()
	})

	return model
}

func newConnectionFormModel() *Model {
	return &Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Database of choice? *").
					Options(
						huh.NewOption(db.Sqlite, "sqlite3"),
						huh.NewOption(db.Oracle, "godror"),
						huh.NewOption(db.Postgres, "pgx"),
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

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.form.Update(msg)
	m.form = updated.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, tea.Sequence(common.SwitchTo(common.ConnectionEstablishmentScreen), func() tea.Msg {
			return common.ConnectionInfo{
				Driver: m.form.GetString(selectedDriver),
				Url:    m.form.GetString(databaseUrl),
				Name:   m.form.GetString(dbName),
			}
		})
	}

	return m, cmd
}

func (m *Model) View() tea.View {
	return tea.NewView(m.form.View())
}
