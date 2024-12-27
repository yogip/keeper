package tui

import (
	"fmt"
	"keeper/internal/infra/tui/views"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Viewer interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
}

type Client interface {
	views.ClientIAM
	views.ClientApp
}

type updateErrMsg string

func updatErrListCmd() tea.Msg {
	return updateErrMsg("")
}

type Model struct {
	app            Client
	screen         views.ScreenType
	viewLogin      *views.LoginView
	viewSignUp     *views.SignUpView
	viewSecretList *views.SecretListView
	viewSecretView *views.SecretView
	viewNewSecret  *views.CreateSecretView
	viewNewPwd     *views.UpsertPwdView
	viewNewNote    *views.UpsertNoteView
	viewNewFile    *views.UpsertPwdView
	viewNewCard    *views.UpsertCardView
	activeView     Viewer
	errors         []*views.ErrorMsg
}

func InitModel(app Client) Model {
	l := views.NewLoginView(app)
	return Model{
		app:            app,
		screen:         views.ScreenLogin,
		viewLogin:      l,
		viewSignUp:     views.NewSignUpView(app),
		viewSecretList: views.NewSecretList(app),
		viewSecretView: views.NewSecretView(app),
		viewNewSecret:  views.NewCreateSecretView(app),
		viewNewPwd:     views.NewUpsertPwdView(app),
		viewNewNote:    views.NewUpsertNoteView(app),
		viewNewFile:    views.NewUpsertPwdView(app), // todo
		viewNewCard:    views.NewUpsertCardView(app),
		activeView:     l,
		errors:         make([]*views.ErrorMsg, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case views.ScreenTypeMsg:
		return m, m.changeViewer(msg)
	case views.ErrorMsg:
		log.Println("Got error", msg.Error())
		m.errors = append(m.errors, &msg)
	}

	errCmd := m.tickErrors()
	cmd := m.activeView.Update(msg)

	return m, tea.Batch(cmd, errCmd)
}

func (m Model) View() string {
	v := m.activeView.View()
	if len(m.errors) == 0 {
		return v
	}

	// Render errors
	var b strings.Builder
	b.WriteString(v)
	b.WriteRune('\n')

	b.WriteString("----------------")
	b.WriteRune('\n')
	for _, err := range m.errors {
		b.WriteString(views.ErrorStyle.Render(err.Error()))
		b.WriteString(fmt.Sprintf(" (%d)", err.HideAterSec()))
		b.WriteString("\n----------------\n")
	}

	return b.String()
}

func (m *Model) tickErrors() tea.Cmd {
	if len(m.errors) == 0 {
		return nil
	}
	errors := make([]*views.ErrorMsg, 0, len(m.errors))
	for _, err := range m.errors {
		if err.IsShowing() {
			errors = append(errors, err)
		}
	}
	m.errors = errors
	return updatErrListCmd
}

func (m *Model) changeViewer(msg views.ScreenTypeMsg) tea.Cmd {
	// log.Printf("Main View. Change Screent to %s from %s\n", string(msg.Screen), string(m.screen))
	m.screen = msg.Screen
	switch msg.Screen {
	// Login
	case views.ScreenLogin:
		m.activeView = m.viewLogin
	// Sign Up
	case views.ScreenSignUp:
		m.activeView = m.viewSignUp

	// List View
	case views.ScreenSecretList:
		m.activeView = m.viewSecretList
		return m.viewSecretList.Init()
	// Secret View
	case views.ScreenSecretView:
		m.activeView = m.viewSecretView
		return m.viewSecretView.Init(*msg.SecretID)

	// New Secret (select secret type)
	case views.ScreenNewSecret:
		m.activeView = m.viewNewSecret

	// New Password
	case views.ScreenUpsertPassword:
		m.activeView = m.viewNewPwd
		return m.viewNewPwd.Init(msg.SecretID)
	// New Note
	case views.ScreenUpsertNote:
		m.activeView = m.viewNewNote
		return m.viewNewNote.Init(msg.SecretID)
	// New card
	case views.ScreenUpsertCard:
		m.activeView = m.viewNewCard
		return m.viewNewCard.Init(msg.SecretID)
	// New File
	case views.ScreenUpsertFile:
		m.activeView = m.viewNewFile
	}
	return nil
}
