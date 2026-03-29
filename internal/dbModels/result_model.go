package dbModels

import (
	"qit/internal/common"
	"qit/internal/db"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ResultModel struct {
	table      table.Model
	spinner    spinner.Model
	loading    bool
	err        error
	page       int
	pageSize   int
	totalRows  int
	sortColumn string
	sortAsc    bool
}

func NewResultModel() *ResultModel {
	spinnerModel := spinner.New(spinner.WithSpinner(spinner.Dot))
	tableModel := table.New(table.WithFocused(true))
	return &ResultModel{
		spinner: spinnerModel,
		table:   tableModel,
		loading: true,
	}
}

func (m *ResultModel) EnterParams(params common.AdditionalParam) {}

func (m *ResultModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, db.ExecuteQueryInDatabase(query, isTransactional))
}

func (m *ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case db.QueryResultMsg:
		m.loading = false
		m.err = msg.Err
		m.table.SetColumns(msg.Columns)
		m.table.SetRows(msg.Rows)
		m.table.SetHeight(len(msg.Rows) - 2)
		m.table.SetCursor(0)

	case spinner.TickMsg:
		// allow spinner to animate
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case table.KeyMap:

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 2)

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ResultModel) View() string {
	if m.loading {
		return "Running given query... " + m.spinner.View()
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress q to quit."
	}

	return m.table.View() + "\n" + m.table.Help.View(table.DefaultKeyMap())
}
