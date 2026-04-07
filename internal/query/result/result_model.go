package result

import (
	"query-terminal/internal/common"
	"query-terminal/internal/db"
	"query-terminal/internal/query/editor"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Quit     key.Binding
}

// ShortHelp implements the KeyMap interface.
func (km keyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.LineUp, km.LineDown, km.Quit}
}

// FullHelp implements the KeyMap interface.
func (km keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.LineUp, km.LineDown, km.PageUp, km.PageDown},
		{km.Quit},
	}
}

// DefaultKeyMap returns a default set of keybindings.
func defaultKeyMap() keyMap {
	const spacebar = " "
	return keyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit program"),
		),
	}
}

type Model struct {
	table      table.Model
	spinner    spinner.Model
	keys       keyMap
	waiting    bool
	loading    bool
	err        error
	page       int
	pageSize   int
	totalRows  int
	sortColumn string
	sortAsc    bool
}

func NewResultModel() Model {
	spinnerModel := spinner.New(spinner.WithSpinner(spinner.Dot))
	tableModel := table.New(table.WithFocused(true)) //otherwise we wouldn't be able to scroll through the table
	return Model{
		spinner: spinnerModel,
		table:   tableModel,
		keys:    defaultKeyMap(),
		waiting: true,
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		// allow spinner to animate
		if m.loading || m.waiting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 2)

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

	case editor.UserInputMsg:
		m.waiting = false
		m.loading = true
		cmds = append(cmds, db.ExecuteQueryInDatabase(msg.InputQuery))

	case db.QueryResultMsg:
		m.loading = false
		m.err = msg.Err
		m.table.SetColumns(msg.Columns)
		m.table.SetRows(msg.Rows)
		m.table.SetHeight(11) // 11 since 1 row is the header, and we want to show 10 rows of data
		m.table.SetCursor(0)

	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.waiting {
		return "Waiting for input query... " + m.spinner.View()
	}

	if m.loading {
		return "Running given query... " + m.spinner.View()
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress ctrl+c to quit."
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.table.View(),
		m.table.Help.View(m.keys),
	)
}
