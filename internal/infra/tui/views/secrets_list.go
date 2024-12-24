package views

import (
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strconv"
	"strings"
	"time"

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
	app            ClientApp
}

func NewSecretList(app ClientApp) *SecretListView {
	searchInput := textinput.New()
	searchInput.Cursor.Style = cursorStyle
	searchInput.CharLimit = 50
	searchInput.Placeholder = "Secret name"
	searchInput.Focus()
	searchInput.PromptStyle = focusedStyle
	searchInput.TextStyle = focusedStyle

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
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
		app:            app,
	}
}

func (m *SecretListView) Init() tea.Cmd {
	return m.searchSecrets("")
}

func (m *SecretListView) searchSecrets(secretName string) tea.Cmd {
	return func() tea.Msg {
		l, err := m.app.ListSecrets(secretName)
		if err != nil {
			log.Println("Error loading secrets", err)
			return NewErrorMsg(err, time.Second*30)
		}

		rows := make([]table.Row, 0, len(l.Secrets))
		for _, s := range l.Secrets {
			rows = append(rows, table.Row{fmt.Sprint(s.ID), s.Name, string(s.Type)})
		}
		m.table.SetRows(rows)
		return ""
	}
}

func (m *SecretListView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			if m.focusIndex == m.focusTable && row != nil {
				sid, err := strconv.ParseInt(row[0], 10, 64)
				if err != nil {
					return ErrorCmd(err, time.Second*10)
				}
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretView, SecretID: &sid})
			}
		// Set focus to next input
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			var skip bool
			// If User taps up button while Table is active and cursor > 0
			// then skip focusIndex increment
			if m.focusTable == m.focusIndex && (s == "shift+tab" || s == "up") && m.table.Cursor() > 0 {
				skip = true
			}
			// If User taps down button while Table is active and cursor < len(rows)-1
			// then skip focusIndex increment
			if m.focusTable == m.focusIndex && (s == "tab" || s == "down") && m.table.Cursor() < len(m.table.Rows())-1 {
				skip = true
			}

			// Cycle indexes
			if !skip {
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
				log.Println("Table focused")
				m.table.Focus()
				m.table, cmd = m.table.Update(msg)
			} else {
				// Remove focused state
				log.Println("Table blured")
				m.table.Blur()
			}

			return cmd
		default:
			if m.focusIndex == m.focusSearchInp {
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				return tea.Batch(
					cmd,
					m.searchSecrets(m.searchInput.Value()),
				)
			}
			// Swtich to create secret view
			if msg.String() == "n" {
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenNewSecret})
			}
			// Swtich to edit secret view
			if msg.String() == "e" {
				row := m.table.SelectedRow()
				if row != nil {
					sid, err := strconv.ParseInt(row[0], 10, 64)
					if err != nil {
						return ErrorCmd(err, time.Second*10)
					}
					switch model.SecretType(row[2]) {
					case model.SecretTypePassword:
						return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertPassword, SecretID: &sid})
					case model.SecretTypeNote:
						return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertNote, SecretID: &sid})
					case model.SecretTypeCard:
						return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertCard, SecretID: &sid})
					case model.SecretTypeFile:
						return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertFile, SecretID: &sid})
					}
				}
			}
		}
	}

	// var cmd tea.Cmd
	// m.table, cmd = m.table.Update(msg)
	return nil
}

func (m *SecretListView) View() string {
	var b strings.Builder
	b.WriteString(boldStyle.Render("Search secret:"))
	b.WriteRune('\n')

	// search text input
	b.WriteString(m.searchInput.View())
	b.WriteRune('\n')

	// secrets table
	b.WriteString(baseStyle.Render(m.table.View()))
	b.WriteRune('\n')

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `up` and `down` or `tab` and `shift+tab` to navigate"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `e` to edit secret, `enter` to see secret info "))
	b.WriteString(helpStyle.Render("and `n` to create a new one."))

	return b.String()
}
