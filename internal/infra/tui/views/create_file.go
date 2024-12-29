package views

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type CreateFileView struct {
	focusIndex      int
	focusMax        int
	focusName       int
	focusFilePicker int
	focusNote       int
	focusSubmit     int
	focusCancel     int
	selectedFile    string
	nameInput       textinput.Model
	filePicker      filepicker.Model
	noteInput       textarea.Model

	secretID *int64
	app      ClientApp
}

func NewCreateFileView(app ClientApp) *CreateFileView {
	// Name
	nameInput := textinput.New()
	nameInput.CharLimit = 50
	nameInput.Placeholder = "Secret name"
	nameInput.Focus()
	nameInput.Cursor.Style = cursorStyle
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	// File picker
	fp := filepicker.New()
	fp.Height = 5
	fp.CurrentDirectory, _ = os.UserHomeDir()
	log.Println("Set current directory to", fp.CurrentDirectory)

	// Note
	noteInput := textarea.New()
	noteInput.Cursor.Style = cursorStyle
	noteInput.CharLimit = 50
	noteInput.Placeholder = "Enter a note"
	noteInput.Blur()

	return &CreateFileView{
		focusIndex: 0,
		focusMax:   8,

		focusName:       0,
		focusFilePicker: 1,
		focusNote:       2,
		focusSubmit:     3,
		focusCancel:     4,

		nameInput:  nameInput,
		filePicker: fp,
		noteInput:  noteInput,

		app: app,
	}
}

func (m *CreateFileView) Init() tea.Cmd {
	log.Println("CreateFileView::Init filePicker was initialized")
	return m.filePicker.Init()
}

func (m *CreateFileView) createSecretCmd(name, filePath, note string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(filePath)
		if err != nil {
			log.Println("Opening file error", err)
			return NewErrorMsg(err, time.Second*15)
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Println("Closing file error", err)
			}
		}()

		body, err := io.ReadAll(file)
		if err != nil {
			log.Println("Reading file error", err)
			return NewErrorMsg(err, time.Second*15)
		}

		fileParts := strings.Split(filePath, "/")
		fileName := fileParts[len(fileParts)-1]

		_, err = m.app.CreateFileSecret(name, fileName, note, body)
		if err != nil {
			log.Println("call grpc method error", err)
			return NewErrorMsg(err, time.Second*10)
		}
		log.Println("Succesfully upsert secret", name)
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	}
}

func (m *CreateFileView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenNewSecret})
		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// Move msg to the picker until a file isn't selected
			if m.selectedFile == "" && m.focusIndex == m.focusFilePicker {
				didSelect, selectedFile := m.filePicker.DidSelectFile(msg)
				if didSelect {
					m.selectedFile = selectedFile
				}
				break
			}

			// Allow to do a line break for text areas elements
			if s == "enter" && m.focusIndex == m.focusNote {
				break
			}

			if s == "enter" && m.focusIndex == m.focusCancel {
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
			}
			if s == "enter" && m.focusIndex == m.focusSubmit {
				// todo
				cmds := make([]tea.Cmd, 0, 5)
				if m.nameInput.Value() == "" {
					cmds = append(cmds, ErrorCmd(errors.New("Secret Name cannot be empty"), time.Second*15))
				}

				if m.selectedFile == "" {
					cmds = append(cmds, ErrorCmd(errors.New("Select a file"), time.Second*15))
				}

				if len(cmds) > 0 {
					return tea.Batch(cmds...)
				}
				return m.createSecretCmd(
					m.nameInput.Value(),
					m.selectedFile,
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

			// todo File picker Component
			if m.focusIndex == m.focusFilePicker {
				// cmd = m.filePicker.Focus()
			} else {
			}

			// Note Component
			if m.focusIndex == m.focusNote {
				cmd = m.noteInput.Focus()
			} else {
				m.noteInput.Blur()
			}

			return cmd
		}
	}
	return m.updateInputs(msg)
}

func (m *CreateFileView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.noteInput, cmds[1] = m.noteInput.Update(msg)
	m.filePicker, cmds[2] = m.filePicker.Update(msg)
	return tea.Batch(cmds...)
}

func (m *CreateFileView) View() string {
	var b strings.Builder
	action := "Create"
	if m.secretID != nil {
		action = "Update"
	}
	b.WriteString(boldStyle.Render(action + " Card:"))
	b.WriteRune('\n')

	b.WriteString(m.nameInput.View())
	b.WriteRune('\n')
	b.WriteRune('\n')

	if m.selectedFile == "" {
		b.WriteString(boldStyle.Render("Select a file:"))
		b.WriteRune('\n')
		b.WriteRune('\n')
		b.WriteString(m.filePicker.View())
	} else {
		b.WriteString(boldStyle.Render("File: ") + m.selectedFile)
	}
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
