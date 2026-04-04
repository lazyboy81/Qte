package establish

import (
	"query-terminal/internal/common"
	"query-terminal/internal/db"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	driver string
	url    string
	name   string
)

type Model struct {
	spinner   spinner.Model
	loading   bool
	err       error
	connected bool
}

func NewConnectionEstablishmentModel() Model {
	model := spinner.New(spinner.WithSpinner(spinner.Dot))
	return Model{
		spinner: model,
		loading: true,
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {
	driver = params["driver"].(string)
	url = params["url"].(string)
	name = params["name"].(string)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, db.ConnectToDatabase(driver, url, name))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case db.ConnectResultMsg:
		m.loading = false
		m.err = msg.Err
		m.connected = msg.Err == nil

	case spinner.TickMsg:
		// allow spinner to animate
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	}

	if !m.loading && m.connected {
		cmds = append(cmds, common.SwitchTo(common.QueryEditorScreen, nil))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.loading {
		return "Connecting to server... " + m.spinner.View()
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress ctrl+c to quit."
	}

	return "Connected successfully!"
}
