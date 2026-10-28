package modal

import (
	"query-terminal/internal/common"
	"query-terminal/internal/connection/form"
	"query-terminal/internal/connection/picker"
	"sync"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type menuItem struct {
	title, desc string
	screen      common.SwitchScreenMsg
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type Model struct {
	form   *form.Model
	picker *picker.Model
	list   list.Model
}

var model *Model
var once sync.Once

func GetNewConnectionModal() *Model {
	once.Do(func() {
		model = newConnectionModalModel()
	})

	return model
}

func newConnectionModalModel() *Model {
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
	}

	listModel := list.New(items, list.NewDefaultDelegate(), 60, 20)
	listModel.Title = "Connection"
	listModel.SetShowStatusBar(false)

	return &Model{
		form:   form.GetConnectionFormModel(),
		picker: picker.GetConnectionPickerModel(),
		list:   listModel,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
	return tea.NewView(m.list.View())
}
