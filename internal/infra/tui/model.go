package tui

import (
	"fmt"
	"keeper/internal/core/model"
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
	Login(req model.UserRequest) (model.Token, error)
	SignUp(req model.UserRequest) (model.Token, error)
}

type updateErrMsg string

func updatErrListCmd() tea.Msg {
	return updateErrMsg("")
}

type Model struct {
	screen         views.ScreenType
	viewLogin      *views.LoginView
	viewSignUp     *views.SignUpView
	viewSecretList *views.SecretListView
	activeView     Viewer
	errors         []*views.ErrorMsg
}

func InitModel(app Client) Model {
	l := views.NewLoginView(app)
	return Model{
		screen:         views.ScreenLogin,
		viewLogin:      l,
		viewSignUp:     views.NewSignUpView(app),
		viewSecretList: views.NewSecretList(),
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
		log.Println("Login:", msg.Token)
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

	var b strings.Builder
	b.WriteString("----------------")
	b.WriteRune('\n')
	for _, err := range m.errors {
		b.WriteString(views.ErrorStyle.Render(err.Error()))
		b.WriteString(fmt.Sprintf(" (%d)", err.HideAterSec()))
		b.WriteString("\n----------------\n")
	}
	b.WriteRune('\n')

	b.WriteString(v)
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
	m.screen = active
	switch active {
	case views.ScreenLogin:
		m.activeView = m.viewLogin
	case views.ScreenSignUp:
		m.activeView = m.viewSignUp
	case views.ScreenSecretList:
		m.activeView = m.viewSecretList
	}
}
