package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/styles"
)

// TextInputComponent handles text input with optional placeholder
type TextInputComponent struct {
	textInput   textinput.Model
	label       string
	placeholder string
	helpText    string
	value       string
	submitted   bool
	cancelled   bool
	optional    bool
}

func NewTextInputComponent(label, placeholder, helpText string, optional bool) *TextInputComponent {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 60

	return &TextInputComponent{
		textInput:   ti,
		label:       label,
		placeholder: placeholder,
		helpText:    helpText,
		optional:    optional,
	}
}

func (t *TextInputComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TextInputComponent) Update(msg tea.Msg) (*TextInputComponent, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			t.cancelled = true
			return t, tea.Quit

		case "enter":
			// Submit the value (can be empty if optional)
			t.value = strings.TrimSpace(t.textInput.Value())
			t.submitted = true
			return t, tea.Quit
		}
	}

	t.textInput, cmd = t.textInput.Update(msg)
	return t, cmd
}

func (t *TextInputComponent) View() string {
	var b strings.Builder

	// Label
	b.WriteString(styles.ActiveStyle.Render(t.label))
	b.WriteString("\n\n")

	// Input field
	b.WriteString(t.textInput.View())
	b.WriteString("\n")

	// Help text
	b.WriteString("\n")
	if t.helpText != "" {
		b.WriteString(styles.FaintStyle.Render(t.helpText))
		b.WriteString("\n")
	}
	if t.optional {
		b.WriteString(styles.FaintStyle.Render("Press Enter to confirm (leave empty to skip), Esc to cancel"))
	} else {
		b.WriteString(styles.FaintStyle.Render("Press Enter to confirm, Esc to cancel"))
	}

	return b.String()
}

func (t *TextInputComponent) Submitted() bool {
	return t.submitted
}

func (t *TextInputComponent) Cancelled() bool {
	return t.cancelled
}

func (t *TextInputComponent) Value() string {
	return t.value
}
