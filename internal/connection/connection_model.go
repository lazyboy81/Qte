package connection

import (
	"TQT/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	url      string
	username string
	pass     string
)

type ConnectionModel struct {
	form *huh.Form
}

func NewConnectionModel() ConnectionModel {
	return ConnectionModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("DB URL").Value(&url),
				huh.NewInput().Title("Username").Value(&username),
				huh.NewInput().Title("Password").Value(&pass),
			),
		),
	}
}

func (m ConnectionModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	update, cmd := m.form.Update(msg)

	if m.form.State == huh.StateCompleted {
		return m, func() tea.Msg { return common.ScreenHome }
	}

	return update, cmd
}

func (m ConnectionModel) View() string {
	return m.form.View()
}
