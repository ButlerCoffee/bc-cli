package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/tui/prompts"
)

// BrewOption represents a brewing method option
type BrewOption struct {
	Value       string
	Display     string
	Description string
	ShowGrind   bool
}

func (b BrewOption) Label() string {
	if b.ShowGrind && b.Description != "" {
		return fmt.Sprintf("%s (%s)", b.Display, b.Description)
	}
	return b.Display
}

func (b BrewOption) Details() string {
	return ""
}

// BrewSelectorModel composes duck + select for brewing method selection
type BrewSelectorModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewBrewSelectorModel(grindType string) BrewSelectorModel {
	showGrind := (grindType == "ground")

	options := []BrewOption{
		{"espresso", "Espresso", "very fine grind", showGrind},
		{"moka", "Moka Pot", "fine-medium grind", showGrind},
		{"v60", "V60 Pour Over", "medium grind", showGrind},
		{"french_press", "French Press", "coarse grind", showGrind},
		{"pour_over", "Pour Over", "medium grind", showGrind},
		{"drip", "Drip Coffee", "medium grind", showGrind},
		{"cold_brew", "Cold Brew", "extra coarse grind", showGrind},
	}

	items := make([]components.SelectItem, len(options))
	for i, opt := range options {
		items[i] = opt
	}

	title := "Select your brewing method"
	if grindType == "ground" {
		title = "Select your brewing method (grind size shown)"
	}

	return BrewSelectorModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent(title, items),
	}
}

func (m BrewSelectorModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m BrewSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update duck (handles tick messages)
	var duckCmd tea.Cmd
	m.duck, duckCmd = m.duck.Update(msg)
	if duckCmd != nil {
		cmds = append(cmds, duckCmd)
	}

	// Update selector (handles key messages)
	var selectCmd tea.Cmd
	m.selector, selectCmd = m.selector.Update(msg)
	if selectCmd != nil {
		cmds = append(cmds, selectCmd)
	}

	// Trigger duck action on selection
	if m.selector.Selected() {
		m.duck.TriggerAction()
	}

	return m, tea.Batch(cmds...)
}

func (m BrewSelectorModel) View() string {
	return m.duck.View() + m.selector.View()
}

// BrewingMethodResult contains the selected brewing method and optional notes
type BrewingMethodResult struct {
	Method string
	Notes  string
}

// SelectBrewingMethod shows the brew selector and returns the selected brewing method and notes
func SelectBrewingMethod(grindType string) (*BrewingMethodResult, error) {
	p := tea.NewProgram(NewBrewSelectorModel(grindType))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(BrewSelectorModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	brewOpt := selectedItem.(BrewOption)

	// Prompt for optional notes
	notes, err := prompts.PromptText(
		"Add custom notes for this preference (optional)",
		"i.e. I like italian style espresso, I love fruity coffee, etc.",
		"Share any preferences to help us customize your coffee experience",
		true,
	)
	if err != nil && err != prompts.ErrUserCancelled {
		return nil, err
	}
	// If user cancelled, notes will be empty string which is fine

	return &BrewingMethodResult{
		Method: brewOpt.Value,
		Notes:  notes,
	}, nil
}
