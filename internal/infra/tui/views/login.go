package views

import (
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginMsg struct {
	Token model.Token
}

type LoginView struct {
	focusIndex  int
	focusMax    int
	focusSubmit int
	focusSignUp int
	inputs      []textinput.Model
	iam         IAMClient
}

func NewLoginView(iam IAMClient) *LoginView {
	inputs := make([]textinput.Model, 2)

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
		}

		inputs[i] = t
	}
	return &LoginView{
		focusIndex:  0,
		focusMax:    len(inputs) + 1,
		focusSubmit: len(inputs),
		focusSignUp: len(inputs) + 1,
		inputs:      inputs,
		iam:         iam,
	}
}

func (m *LoginView) loginCmd(user string, password string) tea.Cmd {
	return func() tea.Msg {
		token, err := m.iam.Login(model.UserRequest{Login: user, Password: password})
		if err != nil {
			log.Println("Login failed.", err)
			return NewErrorMsg(err, time.Second*10)
		}
		return LoginMsg{Token: token}
	}
}

func (m *LoginView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused
			if s == "enter" && m.focusIndex == m.focusSubmit {
				return m.loginCmd(m.inputs[0].Value(), m.inputs[1].Value())
			}
			// Did the user press enter while the sign up button was focused
			if s == "enter" && m.focusIndex == m.focusSignUp {
				return changeScreenCmd(ScreenSignUp)
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

func (m *LoginView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *LoginView) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Enter your credentials:"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	// text inputs
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	// submit button
	b.WriteString("\n\n")
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "%s", *button)

	// sign up button
	b.WriteString("\n")
	signUpBtn := blurredStyle.Render("[ Sign Up ]")
	if m.focusIndex == len(m.inputs)+1 {
		signUpBtn = fmt.Sprintf("[ %s ]", focusedStyle.Render("Sign Up"))
	}
	fmt.Fprintf(&b, "%s", signUpBtn)

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `up` and `down` or `tab` and `shift+tab` to navigate."))

	return b.String()
}
