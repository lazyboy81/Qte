package establish

import (
	"query-terminal/internal/common"
	"query-terminal/internal/db"
	"sync"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	spinner   spinner.Model
	loading   bool
	err       error
	connected bool
}

var model *Model
var once sync.Once

func GetConnectionEstablishmentModel() *Model {
	once.Do(func() {
		model = newConnectionEstablishmentModel()
	})

	return model
}

func newConnectionEstablishmentModel() *Model {
	model := spinner.New(spinner.WithSpinner(spinner.Dot))
	return &Model{
		spinner: model,
		loading: true,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case common.ConnectionInfo:
		cmds = append(cmds, db.ConnectToDatabase(msg.Driver, msg.Url, msg.Name, 0))

	case *db.ConnectResultMsg:
		m.loading = false
		m.err = msg.Err
		m.connected = msg.Err == nil
		if !m.loading && m.connected {
			cmds = append(cmds, common.SwitchTo(common.QueryScreen), func() tea.Msg {
				return msg.Conn
			})
		}

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

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	return tea.NewView("Connecting to server... " + m.spinner.View())
}
