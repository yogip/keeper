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

type UpsertNoteView struct {
	focusIndex  int
	focusMax    int
	focusName   int
	focusText   int
	focusNote   int
	focusSubmit int
	focusCancel int

	nameInput textinput.Model
	textInput textarea.Model
	noteInput textarea.Model

	secretID *int64
	app      ClientApp
}

func NewUpsertNoteView(app ClientApp) *UpsertNoteView {
	// Name
	nameInput := textinput.New()
	nameInput.Cursor.Style = cursorStyle
	nameInput.CharLimit = 50
	nameInput.Placeholder = "Secret name"
	nameInput.Focus()
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	// Text
	textInput := textarea.New()
	textInput.Cursor.Style = cursorStyle
	textInput.Placeholder = "Enter a text"
	textInput.Blur()

	// Note
	noteInput := textarea.New()
	noteInput.Cursor.Style = cursorStyle
	noteInput.CharLimit = 50
	noteInput.Placeholder = "Enter a note"
	noteInput.Blur()

	return &UpsertNoteView{
		focusIndex: 0,
		focusMax:   4, // 0 - nameInput, 1 - text, 2 - note, 3 - Submit, 4 - cancel

		focusName:   0,
		focusText:   1,
		focusNote:   2,
		focusSubmit: 3,
		focusCancel: 4,

		nameInput: nameInput,
		textInput: textInput,
		noteInput: noteInput,

		app: app,
	}
}

func (m *UpsertNoteView) Init(secretID *int64) tea.Cmd {
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
		note, err := secret.AsNote()
		if err != nil {
			log.Println("Could not open secret for editing", err)
			return NewErrorMsg(err, time.Second*30)
		}

		m.secretID = &note.ID
		m.nameInput.SetValue(note.Name)
		m.textInput.SetValue(note.Text)
		m.noteInput.SetValue(note.Note)
		log.Println("Secret data loaded:", note.ID, note.Name)
		return ""
	}
}

func (m *UpsertNoteView) upsertSecretCmd(name, text, note string) tea.Cmd {
	return func() tea.Msg {
		s := model.NewNote(0, name, text, note)
		payload, err := s.GetPayload()
		if err != nil {
			log.Println("creating note payload error", err)
			return NewErrorMsg(err, time.Second*10)
		}

		if m.secretID == nil {
			_, err = m.app.CreateSecret(model.SecretTypeNote, name, note, payload)
		} else {
			_, err = m.app.UpdateSecret(*m.secretID, model.SecretTypeNote, name, note, payload)
		}
		if err != nil {
			log.Println("call grpc method error", err)
			return NewErrorMsg(err, time.Second*10)
		}
		log.Println("Succesfully upsert secret", name)
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	}
}

func (m *UpsertNoteView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenNewSecret})
		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// Allow to do a line break for text areas elements
			if s == "enter" && (m.focusIndex == m.focusNote || m.focusIndex == m.focusText) {
				return m.updateInputs(msg)
			}

			if s == "enter" && m.focusIndex == m.focusCancel {
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
			}
			if s == "enter" && m.focusIndex == m.focusSubmit {
				if m.nameInput.Value() == "" || m.textInput.Value() == "" {
					return ErrorCmd(errors.New("Secret Name and Text cannot be empty"), time.Second*5)
				}
				return m.upsertSecretCmd(
					m.nameInput.Value(),
					m.textInput.Value(),
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
			if m.focusIndex == m.focusText {
				cmd = m.textInput.Focus()
				// m.textInput.FocusedStyle = focusedStyle
				// m.textInput.TextStyle = focusedStyle
			} else {
				m.textInput.Blur()
				// m.textInput.PromptStyle = noStyle
				// m.textInput.TextStyle = noStyle
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

func (m *UpsertNoteView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.textInput, cmds[1] = m.textInput.Update(msg)
	m.noteInput, cmds[2] = m.noteInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m *UpsertNoteView) View() string {
	var b strings.Builder
	action := "Create"
	if m.secretID != nil {
		action = "Update"
	}
	b.WriteString(boldStyle.Render(action + " Note:"))
	b.WriteRune('\n')

	b.WriteString(m.nameInput.View())
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(m.textInput.View())
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
