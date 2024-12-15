package views

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type CreatePwdView struct {
	focusIndex  int
	focusMax    int
	focusName   int
	focusLogin  int
	focusPwd    int
	focusNote   int
	focusSubmit int
	focusCancel int

	nameInput  *textinput.Model
	loginInput *textinput.Model
	pwdInput   *textinput.Model
	noteInput  *textinput.Model
	inputs     []*textinput.Model

	app ClientApp
}

func NewCreatePwdView(app ClientApp) *CreatePwdView {
	// Name
	nameInput := textinput.New()
	nameInput.Cursor.Style = cursorStyle
	nameInput.CharLimit = 50
	nameInput.Placeholder = "Secret name"
	nameInput.Focus()
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	// Login
	loginInput := textinput.New()
	loginInput.Cursor.Style = cursorStyle
	loginInput.CharLimit = 50
	loginInput.Placeholder = "Login"
	loginInput.Blur()
	loginInput.PromptStyle = focusedStyle
	loginInput.TextStyle = focusedStyle

	// Password
	pwdInput := textinput.New()
	pwdInput.Cursor.Style = cursorStyle
	pwdInput.CharLimit = 50
	pwdInput.Placeholder = "Password"
	pwdInput.Blur()
	pwdInput.PromptStyle = focusedStyle
	pwdInput.TextStyle = focusedStyle

	// Note
	noteInput := textinput.New()
	noteInput.Cursor.Style = cursorStyle
	noteInput.CharLimit = 50
	noteInput.Placeholder = "Note"
	noteInput.Blur()
	noteInput.PromptStyle = focusedStyle
	noteInput.TextStyle = focusedStyle

	return &CreatePwdView{
		focusIndex: 0,
		focusMax:   5, // 0 - nameInput, 1 - login, 2 - pwd, 3 - note, 4 - Submit, 5 - cancel

		focusName:   0,
		focusLogin:  1,
		focusPwd:    2,
		focusNote:   3,
		focusSubmit: 4,
		focusCancel: 5,

		nameInput:  &nameInput,
		loginInput: &loginInput,
		pwdInput:   &pwdInput,
		noteInput:  &noteInput,
		inputs:     []*textinput.Model{&nameInput, &loginInput, &pwdInput, &noteInput},

		app: app,
	}
}

func (m *CreatePwdView) Init() tea.Cmd {
	return nil
}

func (m *CreatePwdView) Update(msg tea.Msg) tea.Cmd {
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
				if m.nameInput.Value() == "" || m.pwdInput.Value() == "" {
					return ErrorCmd(errors.New("Secret Name and Password cannot be empty"), time.Second*5)
				}
				// todo calse grpc
			}

			// Cycle indexes
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

			var cmd tea.Cmd
			cmd = m.handleFocus(m.focusIndex)

			return cmd
		}

		return m.updateInputs(msg)
	}
	return nil
}

func (m *CreatePwdView) handleFocus(index int) tea.Cmd {
	var cmd tea.Cmd
	for i, inp := range m.inputs {
		log.Printf("createPwd handleFocus i: %d, index: %d, ", i, index)
		if i != index {
			blurInput(inp)
		}
		if i == index {
			cmd = focusInput(inp)
		}
	}
	return cmd
}

func (m *CreatePwdView) updateInputs(msg tea.Msg) tea.Cmd {
	if m.focusIndex >= len(m.inputs) {
		return nil
	}
	inp, cmd := m.inputs[m.focusIndex].Update(msg)
	m.inputs[m.focusIndex] = &inp

	return cmd
}

func (m *CreatePwdView) View() string {
	var b strings.Builder
	b.WriteString(boldStyle.Render("Create password:"))
	b.WriteRune('\n')

	b.WriteString(m.nameInput.View())
	b.WriteRune('\n')

	b.WriteString(m.loginInput.View())
	b.WriteRune('\n')

	b.WriteString(m.pwdInput.View())
	b.WriteRune('\n')

	b.WriteString(m.noteInput.View())
	b.WriteRune('\n')

	// submit button
	b.WriteString("\n")
	button := blurredStyle.Render("[ Create ]")
	if m.focusIndex == m.focusSubmit {
		button = fmt.Sprintf("[ %s ]", focusedStyle.Render("Create"))
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
