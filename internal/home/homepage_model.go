package home

import (
	"query-terminal/internal/common"
	"sync"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	listStyle = lipgloss.NewStyle().Padding(1, 2)
)

type menuItem struct {
	title, desc string
	screen      common.SwitchScreenMsg
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type Model struct {
	help help.Model
	list list.Model
}

var model *Model
var once sync.Once

func GetHomePage() *Model {
	once.Do(func() {
		model = newHomePage()
	})

	return model
}

func newHomePage() *Model {
	items := []list.Item{
		menuItem{
			title:  "New Connection",
			desc:   "Connect to a database",
			screen: common.ConnectionFormScreen,
		},
		menuItem{
			title:  "Existing Connection",
			desc:   "Connect to an existing database",
			screen: common.ConnectionPickerScreen,
		},
		menuItem{
			title:  "Info",
			desc:   "Information about this application",
			screen: common.InfoScreen,
		},
	}

	listModel := list.New(items, list.NewDefaultDelegate(), 60, 20)
	listModel.Title = "Options"
	listModel.SetShowStatusBar(false)

	return &Model{
		help: help.New(),
		list: listModel,
	}
}

// Init we do any initialization that we want to do in this method (sort of like the love.load() function).
// since we are simply getting user input here which are db url, username and password we don't need anything
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update we switch on the msg type that was given to us and match it
// against the key mappings that we have defined
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// Get selected item and switch to its screen
			selected, ok := m.list.SelectedItem().(menuItem)
			if ok {
				return m, common.SwitchTo(selected.screen)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() tea.View {
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Center,
			listStyle.Render(m.list.View()),
		),
	)
}
