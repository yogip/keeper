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
	viewNewSecret  *views.CreateSecretView
	viewNewPwd     *views.CreatePwdView
	viewNewNote    *views.CreatePwdView
	viewNewFile    *views.CreatePwdView
	viewNewCard    *views.CreatePwdView
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
		viewNewSecret:  views.NewCreateSecretView(app),
		viewNewPwd:     views.NewCreatePwdView(app),
		viewNewNote:    views.NewCreatePwdView(app), // todo
		viewNewFile:    views.NewCreatePwdView(app), // todo
		viewNewCard:    views.NewCreatePwdView(app), // todo
		activeView:     l,
		errors:         make([]*views.ErrorMsg, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case views.ScreenType:
		m.changeViewer(msg)
	case views.LoginMsg:
		m.changeViewer(views.ScreenSecretList)
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

func (m *Model) changeViewer(active views.ScreenType) {
	log.Printf("Main View. Change Screent to %s from %s\n", string(active), string(m.screen))
	m.screen = active
	switch active {
	// Login
	case views.ScreenLogin:
		m.activeView = m.viewLogin
	// Sign Up
	case views.ScreenSignUp:
		m.activeView = m.viewSignUp
	// List View
	case views.ScreenSecretList:
		m.activeView = m.viewSecretList
		m.viewSecretList.Init()
	// New Secret (select secret type)
	case views.ScreenNewSecret:
		m.activeView = m.viewNewSecret

	// New Password
	case views.ScreenNewPassword:
		m.activeView = m.viewNewPwd
	// New Note
	case views.ScreenNewNote:
		m.activeView = m.viewNewNote
	// New card
	case views.ScreenNewCard:
		m.activeView = m.viewNewCard
	// New File
	case views.ScreenNewFile:
		m.activeView = m.viewNewFile
	}
}
