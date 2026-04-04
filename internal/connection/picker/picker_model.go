package picker

import (
	"query-terminal/internal/common"
	"query-terminal/internal/store"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle = lipgloss.NewStyle().Padding(1, 2)
)

type menuItem struct {
	name, url, driver string
}

func (i menuItem) Title() string       { return i.name }
func (i menuItem) Description() string { return i.url }
func (i menuItem) FilterValue() string { return i.name }

type Model struct {
	spinner spinner.Model
	list    *list.Model
	loading bool
	err     error
}

func NewConnectionPickerModel() Model {
	model := spinner.New(spinner.WithSpinner(spinner.Dot))
	return Model{
		spinner: model,
		loading: true,
	}
}

func (m Model) EnterParams(params common.AdditionalParam) {
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, store.ConnectionsCfg().LoadAll())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case store.ConnectionsLoadedMsg:
		m.loading = false
		m.err = msg.Err
		if msg.Err == nil {
			m.list = createListModel(msg.Conns)
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
		if msg.String() == "enter" {
			// Get selected item and switch to its screen
			selected, ok := m.list.SelectedItem().(menuItem)
			if ok {
				return m, common.SwitchTo(common.ConnectionEstablishmentScreen,
					common.AdditionalParam{"driver": selected.driver, "url": selected.url, "name": selected.name},
				)
			}
		}

	}

	if m.list != nil {
		newList, cmd := m.list.Update(msg)
		m.list = &newList
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func createListModel(views []store.ConnectionView) *list.Model {
	items := make([]list.Item, len(views))

	for idx, view := range views {
		items[idx] = menuItem{
			name:   view.Name,
			url:    view.Url,
			driver: view.Driver,
		}
	}

	model := list.New(items, list.NewDefaultDelegate(), 60, 20)
	model.Title = "Connections"
	model.SetShowStatusBar(false)
	return &model
}

func (m Model) View() string {
	if m.loading {
		return "Loading existing connections... " + m.spinner.View()
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress ctrl+c to quit."
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		listStyle.Render(m.list.View()),
	)
}
