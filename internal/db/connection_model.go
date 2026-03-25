package db

import (
	"qit/internal/common"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// messages
type connectResultMsg struct {
	err error
}

// simulate connecting to a foreign server
func connectToDatabase() tea.Cmd {
	return func() tea.Msg {
		// simulate network delay
		_, err := connect(selectedProvider, databaseUrl)

		// return nil on success, or an error
		return connectResultMsg{err: err}
	}
}

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

func (m ConnectionModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, connectToDatabase())
}

func (m ConnectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case connectResultMsg:
		m.loading = false
		m.err = msg.err
		m.connected = msg.err == nil
		return m, nil

	case spinner.TickMsg:
		// allow spinner to animate
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if key.Matches(msg, common.QuitKey) {
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

	return "Connected successfully! 🎉\nPress q to quit."
}
