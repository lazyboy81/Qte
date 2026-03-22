package home

import (
	"TQT/internal/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func NewHomePage() Model {
	items := []list.Item{
		menuItem{
			title:  "New Connection",
			desc:   "Connect to a database",
			screen: common.ScreenConnection,
		},
		menuItem{
			title:  "Info",
			desc:   "Information about this application",
			screen: common.ScreenInfo,
		},
	}

	listModel := list.New(items, list.NewDefaultDelegate(), 60, 20)
	listModel.Title = "Options"
	listModel.SetShowStatusBar(false)

	return Model{
		help: help.New(),
		list: listModel,
	}
}

// Init we do any initialization that we want to do in this method (sort of like the love.load() function).
// since we are simply getting user input here which are db url, username and password we don't need anything
func (m Model) Init() tea.Cmd {
	return nil
}

// Update we switch on the msg type that was given to us and match it
// against the key mappings that we have defined
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.GetBinding("enter")):
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

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		listStyle.Render(m.list.View()),
	)
}
