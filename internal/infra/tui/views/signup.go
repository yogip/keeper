package views

import (
	"errors"
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SignUpView struct {
	focusIndex int
	inputs     []textinput.Model
	iam        ClientIAM
}

func NewSignUpView(iam ClientIAM) *SignUpView {
	inputs := make([]textinput.Model, 3)

	for i := range inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 2:
			t.Placeholder = "Repeat Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}
	return &SignUpView{
		iam:        iam,
		focusIndex: 0,
		inputs:     inputs,
	}
}

func (m *SignUpView) signUpCmd(user string, password string) tea.Cmd {
	return func() tea.Msg {
		err := m.iam.SignUp(model.UserRequest{Login: user, Password: password})
		if err != nil {
			log.Println("Login failed.", err)
			return NewErrorMsg(err, time.Second*10)
		}
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
	}
}

func (m *SignUpView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "esc":
			return tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused
			if s == "enter" && m.focusIndex == len(m.inputs) && m.inputs[1].Value() == m.inputs[2].Value() {
				return m.signUpCmd(m.inputs[0].Value(), m.inputs[1].Value())
			}
			if s == "enter" && m.focusIndex == len(m.inputs) && m.inputs[1].Value() != m.inputs[2].Value() {
				return ErrorCmd(errors.New("Password isn't match"), time.Second*5)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return cmd
}

func (m *SignUpView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *SignUpView) View() string {
	var b strings.Builder
	b.WriteString(boldStyle.Render("Create new account:"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("Use `up` and `down` or `tab` and `shift+tab` to navigate."))

	return b.String()
}
