package views

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type SecretListView struct {
	focusIndex     int
	focusMax       int
	focusSearchInp int
	focusTable     int
	searchInput    textinput.Model
	table          table.Model
}

func NewSecretList() *SecretListView {
	searchInput := textinput.New()
	searchInput.Cursor.Style = cursorStyle
	searchInput.CharLimit = 50
	searchInput.Placeholder = "Secret name"
	searchInput.Focus()
	searchInput.PromptStyle = focusedStyle
	searchInput.TextStyle = focusedStyle

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Type", Width: 10},
	}

	rows := []table.Row{
		{"1", "Name", "Password"},
		{"2", "Name 2", "Password"},
		{"3", "Name 24", "Password"},
		{"4", "Name 22", "Password"},
		{"6", "Name 21", "Password"},
		{"7", "Name 2e", "Password"},
		{"8", "Name 2", "Password"},
		{"9", "Namef2", "Password"},
		{"10", "Namef2", "Password"},
		{"11", "Namde2", "Password"},
		{"12", "Name12", "Password"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	return &SecretListView{
		focusIndex:     0,
		focusMax:       1, // 0 - searchInput, 1 - table
		focusSearchInp: 0,
		focusTable:     1,
		searchInput:    searchInput,
		table:          t,
	}
}

func (m *SecretListView) Update(msg tea.Msg) tea.Cmd {
	log.Println("Update SecretListView", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("Update KeyMsg", msg.String())

		switch msg.String() {
		case "ctrl+c", "esc":
			return tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// if s == "enter" && m.focusIndex == m.focusSubmit {
			// 	return changeScreenCmd(ScreenLogin)
			// }

			// Cycle indexes
			if s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.focusMax {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.focusMax
			}

			var cmd tea.Cmd
			if m.focusIndex == m.focusSearchInp {
				// Set focused state
				cmd = m.searchInput.Focus()
				m.searchInput.PromptStyle = focusedStyle
				m.searchInput.TextStyle = focusedStyle
			} else {
				// Remove focused state
				m.searchInput.Blur()
				m.searchInput.PromptStyle = noStyle
				m.searchInput.TextStyle = noStyle
			}

			if m.focusIndex == m.focusTable {
				// Set focused state
				m.table.Focus()
			} else {
				// Remove focused state
				m.table.Blur()
			}

			return cmd
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return cmd
}

func (m *SecretListView) updateInputs(msg tea.Msg) tea.Cmd {
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	cmds := make([]tea.Cmd, 2)

	m.searchInput, cmds[0] = m.searchInput.Update(msg)
	m.table, cmds[1] = m.table.Update(msg)

	return tea.Batch(cmds...)
}

func (m *SecretListView) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Search secret:"))
	b.WriteRune('\n')

	// search text input
	b.WriteString(m.searchInput.View())
	b.WriteRune('\n')

	// secrets table
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteRune('\n')
	b.WriteString(m.table.HelpView())
	b.WriteRune('\n')

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `tab` to switch between secret name input and table"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `e` to edit secret, `enter` to see secret info "))
	b.WriteString(helpStyle.Render("and `n` to create a new secret."))

	return b.String()
}
