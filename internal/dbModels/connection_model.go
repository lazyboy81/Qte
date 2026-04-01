package dbModels

import (
	"qte/internal/common"
	"qte/internal/db"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	provider string
	url      string
)

type ConnectionModel struct {
	spinner   spinner.Model
	loading   bool
	err       error
	connected bool
}

func NewConnectionModel() ConnectionModel {
	model := spinner.New(spinner.WithSpinner(spinner.Dot))
	return ConnectionModel{
		spinner: model,
		loading: true,
	}
}

func (m ConnectionModel) EnterParams(params common.AdditionalParam) {
	provider = params["provider"].(string)
	url = params["url"].(string)
}

func (m ConnectionModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, db.ConnectToDatabase(provider, url))
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case db.ConnectResultMsg:
		m.loading = false
		m.err = msg.Err
		m.connected = msg.Err == nil
		cmds = append(cmds, common.SwitchTo(common.QueryScreen, nil))

	case spinner.TickMsg:
		// allow spinner to animate
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	}

	return m, tea.Batch(cmds...)
}

func (m ConnectionModel) View() string {
	if m.loading {
		return "Connecting to server... " + m.spinner.View()
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress q to quit."
	}

	return "Connected successfully!"
}
