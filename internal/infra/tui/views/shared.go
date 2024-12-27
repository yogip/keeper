package views

import (
	"fmt"
	"keeper/internal/core/model"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boldStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Bold(true)
	regularStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	ErrorStyle          = blurredStyle.Foreground(lipgloss.Color("#ff0000")).Bold(true)
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	listTitleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Bold(true)
	listTitleStyleFocused = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	listItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	listSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("63")) // 170
)

type ClientApp interface {
	CreateSecret(secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error)
	UpdateSecret(id int64, secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error)
	ListSecrets(secretName string) (*model.SecretList, error)
	GetSecret(secretID int64) (*model.Secret, error)
}

type ScreenType string

type ScreenTypeMsg struct {
	Screen   ScreenType
	SecretID *int64
}

const (
	ScreenLogin          ScreenType = "login"
	ScreenSignUp         ScreenType = "signup"
	ScreenSecretList     ScreenType = "secret-list"
	ScreenSecretView     ScreenType = "secret-view"
	ScreenNewSecret      ScreenType = "new-secret"
	ScreenUpsertPassword ScreenType = "upsert-secret-pwd"
	ScreenUpsertNote     ScreenType = "upsert-secret-note"
	ScreenUpsertCard     ScreenType = "upsert-secret-card"
	ScreenUpsertFile     ScreenType = "upsert-secret-file"
)

func changeScreenCmd(screen *ScreenTypeMsg) tea.Cmd {
	return func() tea.Msg {
		s := ScreenTypeMsg{Screen: screen.Screen, SecretID: screen.SecretID}
		return s
	}
}

type ErrorMsg struct {
	Err       error
	hideAfter time.Time
}

func NewErrorMsg(err error, hideAter time.Duration) ErrorMsg {
	return ErrorMsg{
		Err:       err,
		hideAfter: time.Now().Add(hideAter),
	}
}

func (e *ErrorMsg) Error() string { return e.Err.Error() }

func (e *ErrorMsg) HideAterSec() int {
	return int(e.hideAfter.Sub(time.Now()).Seconds())
}

func (e *ErrorMsg) IsShowing() bool {
	return time.Now().Before(e.hideAfter)
}

func ErrorCmd(err error, hideAter time.Duration) tea.Cmd {
	return func() tea.Msg {
		return NewErrorMsg(err, hideAter)
	}
}
