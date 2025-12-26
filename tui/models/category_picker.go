package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
)

// CategoryItem wraps an api.Category for use with SelectComponent
type CategoryItem struct {
	Category api.Category
	IsExit   bool
}

func (c CategoryItem) Label() string {
	if c.IsExit {
		return "← Exit"
	}
	return c.Category.Name
}

func (c CategoryItem) Details() string {
	if c.IsExit {
		return "Return to main menu"
	}
	return c.Category.Description
}

// CategoryPickerModel composes duck + select for category browsing
type CategoryPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewCategoryPickerModel(categories []api.Category) CategoryPickerModel {
	items := make([]components.SelectItem, 0, len(categories)+1)

	// Add categories
	for _, cat := range categories {
		items = append(items, CategoryItem{Category: cat})
	}

	// Add exit
	items = append(items, CategoryItem{IsExit: true})

	return CategoryPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Explore Coffee Knowledge", items),
	}
}

func (m CategoryPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m CategoryPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var duckCmd tea.Cmd
	m.duck, duckCmd = m.duck.Update(msg)
	if duckCmd != nil {
		cmds = append(cmds, duckCmd)
	}

	var selectCmd tea.Cmd
	m.selector, selectCmd = m.selector.Update(msg)
	if selectCmd != nil {
		cmds = append(cmds, selectCmd)
	}

	if m.selector.Selected() {
		m.duck.TriggerAction()
	}

	return m, tea.Batch(cmds...)
}

func (m CategoryPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickCategory returns selected category or error
func PickCategory(categories []api.Category) (*api.Category, error) {
	p := tea.NewProgram(NewCategoryPickerModel(categories))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(CategoryPickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	catItem := selectedItem.(CategoryItem)
	if catItem.IsExit {
		return nil, nil
	}

	return &catItem.Category, nil
}
