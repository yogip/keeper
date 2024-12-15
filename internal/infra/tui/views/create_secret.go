package views

import (
	"errors"
	"fmt"
	"io"
	"keeper/internal/core/model"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// str := fmt.Sprintf("%d. %s", index+1, i)

	fn := listItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(string(i)))
}

type CreateSecretView struct {
	focusIndex      int
	focusMax        int
	focusSecretType int
	listSecretType  list.Model

	// pwd
	focusLoginInput int
	focusPwdInput   int
	loginInput      textinput.Model
	pwdInput        textinput.Model

	// note
	focusNNoteInput int

	// contorls
	focusSubmit int
	focusCancel int

	// note
	noteInput textinput.Model

	secretType model.SecretType
	app        ClientApp
}

func NewCreateSecretView(app ClientApp) *CreateSecretView {
	items := []list.Item{
		item(model.SecretTypePassword),
		item(model.SecretTypeNote),
		item(model.SecretTypeCard),
		item(model.SecretTypeFile),
	}

	l := list.New(items, itemDelegate{}, 25, 7)
	l.Title = "Select secret type?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.Title = listTitleStyle

	return &CreateSecretView{
		focusIndex:      0,
		focusMax:        2, // 0 - secretType, 1 - Next, 2 - cancel
		focusSecretType: 0,

		focusSubmit:    1,
		focusCancel:    2,
		listSecretType: l,
		app:            app,
	}
}

func (m *CreateSecretView) Init() tea.Cmd {
	return nil
}

func (m *CreateSecretView) nextStageView() tea.Cmd {
	switch m.secretType {
	case model.SecretTypePassword:
		return changeScreenCmd(ScreenNewPassword)
	case model.SecretTypeNote:
		return changeScreenCmd(ScreenNewNote)
	case model.SecretTypeCard:
		return changeScreenCmd(ScreenNewCard)
	case model.SecretTypeFile:
		return changeScreenCmd(ScreenNewFile)
	}
	return ErrorCmd(errors.New("Select secret type"), time.Second*5)
}

func (m *CreateSecretView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			if s == "enter" && m.focusIndex == m.focusCancel {
				return changeScreenCmd(ScreenSecretList)
			}
			if s == "enter" && m.focusIndex == m.focusSubmit {
				return m.nextStageView()
			}
			if s == "enter" && m.focusIndex == m.focusSecretType {
				i, ok := m.listSecretType.SelectedItem().(item)
				if ok {
					m.secretType = model.SecretType(i)
					m.focusIndex++
				}
			}
			skip := m.focusIndex == m.focusSecretType
			indx := m.listSecretType.Index()
			if m.focusIndex == m.focusSecretType && s == "up" && indx == 0 {
				skip = false
			}
			if m.focusIndex == m.focusSecretType && s == "down" && indx == len(m.listSecretType.Items())-1 {
				skip = false
			}

			// Cycle indexes
			if !skip {
				if s == "up" || s == "shift+tab" {
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
			if m.focusIndex == m.focusSecretType {
				m.listSecretType.Styles.Title = listTitleStyleFocused
				m.listSecretType, cmd = m.listSecretType.Update(msg)
			} else {
				m.listSecretType.Styles.Title = listTitleStyle
			}

			return cmd
		}

	}
	return nil
}

func (m *CreateSecretView) View() string {
	var b strings.Builder
	b.WriteString(m.listSecretType.View())

	// Secret type
	b.WriteRune('\n')
	b.WriteString(boldStyle.Render(
		fmt.Sprintf("Secret type: %s\n", m.secretType),
	))
	b.WriteRune('\n')

	// submit button
	b.WriteString("\n")
	button := blurredStyle.Render("[ Next ]")
	if m.focusIndex == m.focusSubmit {
		button = fmt.Sprintf("[ %s ]", focusedStyle.Render("Next"))
	}
	fmt.Fprintf(&b, "%s", button)

	// cancel button
	b.WriteString("\n")
	cancelBtn := blurredStyle.Render("[ Cancel ]")
	if m.focusIndex == m.focusCancel {
		cancelBtn = fmt.Sprintf("[ %s ]", focusedStyle.Render("Cancel"))
	}
	fmt.Fprintf(&b, "%s", cancelBtn)

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `up` and `down` or `tab` and `shift+tab` to navigate"))
	b.WriteString("\n")

	return b.String()
}
