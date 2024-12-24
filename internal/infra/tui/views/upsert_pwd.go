package views

import (
	"errors"
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type UpsertPwdView struct {
	focusIndex  int
	focusMax    int
	focusName   int
	focusLogin  int
	focusPwd    int
	focusNote   int
	focusSubmit int
	focusCancel int

	nameInput  textinput.Model
	loginInput textinput.Model
	pwdInput   textinput.Model
	noteInput  textarea.Model

	secretID *int64
	app      ClientApp
}

func NewUpsertPwdView(app ClientApp) *UpsertPwdView {
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
	loginInput.PromptStyle = noStyle
	loginInput.TextStyle = noStyle

	// Password
	pwdInput := textinput.New()
	pwdInput.Cursor.Style = cursorStyle
	pwdInput.CharLimit = 50
	pwdInput.Placeholder = "Password"
	pwdInput.Blur()
	pwdInput.PromptStyle = noStyle
	pwdInput.TextStyle = noStyle
	pwdInput.EchoMode = textinput.EchoPassword
	pwdInput.EchoCharacter = '•'

	// Note
	noteInput := textarea.New()
	noteInput.Cursor.Style = cursorStyle
	noteInput.CharLimit = 50
	noteInput.Placeholder = "Enter a note"
	noteInput.Blur()

	return &UpsertPwdView{
		focusIndex: 0,
		focusMax:   5, // 0 - nameInput, 1 - login, 2 - pwd, 3 - note, 4 - Submit, 5 - cancel

		focusName:   0,
		focusLogin:  1,
		focusPwd:    2,
		focusNote:   3,
		focusSubmit: 4,
		focusCancel: 5,

		nameInput:  nameInput,
		loginInput: loginInput,
		pwdInput:   pwdInput,
		noteInput:  noteInput,

		app: app,
	}
}

func (m *UpsertPwdView) Init(secretID *int64) tea.Cmd {
	return func() tea.Msg {
		if secretID == nil {
			return ""
		}
		log.Println("Init secret edit view:", *secretID)
		secret, err := m.app.GetSecret(*secretID)
		if err != nil {
			log.Println("Error loading secret", err)
			return NewErrorMsg(err, time.Second*30)
		}
		pwd, err := secret.AsPassword()
		if err != nil {
			log.Println("Could not open secret for editing", err)
			return NewErrorMsg(err, time.Second*30)
		}

		m.secretID = &pwd.ID
		m.nameInput.SetValue(pwd.Name)
		m.loginInput.SetValue(pwd.Login)
		m.pwdInput.SetValue(pwd.Password)
		m.noteInput.SetValue(pwd.Note)
		log.Println("Secret data loaded:", pwd.ID, pwd.Name)
		return ""
	}
}

func (m *UpsertPwdView) upsertPwdCmd(name, login, password, note string) tea.Cmd {
	return func() tea.Msg {
		pwd := model.NewPassword(0, name, login, password, note)
		payload, err := pwd.GetPayload()
		if err != nil {
			log.Println("creating password payload error", err)
			return NewErrorMsg(err, time.Second*10)
		}

		if m.secretID == nil {
			_, err = m.app.CreateSecret(model.SecretTypePassword, name, note, payload)
		} else {
			_, err = m.app.UpdateSecret(*m.secretID, model.SecretTypePassword, name, note, payload)
		}
		if err != nil {
			log.Println("call grpc method error", err)
			return NewErrorMsg(err, time.Second*10)
		}
		log.Println("Succesfully upsert secret", name, login)
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	}
}

func (m *UpsertPwdView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenNewSecret})
		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			if s == "enter" && m.focusIndex == m.focusNote {
				return m.updateInputs(msg)
			}

			if s == "enter" && m.focusIndex == m.focusCancel {
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
			}
			if s == "enter" && m.focusIndex == m.focusSubmit {
				if m.nameInput.Value() == "" || m.pwdInput.Value() == "" {
					return ErrorCmd(errors.New("Secret Name and Password cannot be empty"), time.Second*5)
				}
				return m.upsertPwdCmd(
					m.nameInput.Value(),
					m.loginInput.Value(),
					m.pwdInput.Value(),
					m.noteInput.Value(),
				)
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
			// Name Component
			if m.focusIndex == m.focusName {
				cmd = m.nameInput.Focus()
				m.nameInput.PromptStyle = focusedStyle
				m.nameInput.TextStyle = focusedStyle
			} else {
				m.nameInput.Blur()
				m.nameInput.PromptStyle = noStyle
				m.nameInput.TextStyle = noStyle
			}
			// Login Component
			if m.focusIndex == m.focusLogin {
				cmd = m.loginInput.Focus()
				m.loginInput.PromptStyle = focusedStyle
				m.loginInput.TextStyle = focusedStyle
			} else {
				m.loginInput.Blur()
				m.loginInput.PromptStyle = noStyle
				m.loginInput.TextStyle = noStyle
			}
			// Pwd Component
			if m.focusIndex == m.focusPwd {
				cmd = m.pwdInput.Focus()
				m.pwdInput.PromptStyle = focusedStyle
				m.pwdInput.TextStyle = focusedStyle
			} else {
				m.pwdInput.Blur()
				m.pwdInput.PromptStyle = noStyle
				m.pwdInput.TextStyle = noStyle
			}
			// Note Component
			if m.focusIndex == m.focusNote {
				cmd = m.noteInput.Focus()
			} else {
				m.noteInput.Blur()
			}

			return cmd
		}

		return m.updateInputs(msg)
	}
	return nil
}

func (m *UpsertPwdView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 4)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.loginInput, cmds[1] = m.loginInput.Update(msg)
	m.pwdInput, cmds[2] = m.pwdInput.Update(msg)
	m.noteInput, cmds[3] = m.noteInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m *UpsertPwdView) View() string {
	var b strings.Builder
	action := "Create"
	if m.secretID != nil {
		action = "Update"
	}
	b.WriteString(boldStyle.Render(action + " password:"))
	b.WriteRune('\n')

	b.WriteString(m.nameInput.View())
	b.WriteRune('\n')

	b.WriteString(m.loginInput.View())
	b.WriteRune('\n')

	b.WriteString(m.pwdInput.View())
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(m.noteInput.View())
	b.WriteRune('\n')

	// submit button
	b.WriteString("\n")

	button := blurredStyle.Render(fmt.Sprintf("[ %s ]", action))
	if m.focusIndex == m.focusSubmit {
		button = fmt.Sprintf("[ %s ]", focusedStyle.Render(action))
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
