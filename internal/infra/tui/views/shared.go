package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Bold(true)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	ErrorStyle          = blurredStyle.Foreground(lipgloss.Color("#ff0000")).Bold(true)
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type ScreenType string

const (
	ScreenLogin      ScreenType = "login"
	ScreenSignUp     ScreenType = "signup"
	ScreenSecretList ScreenType = "secret-list"
)

func changeScreenCmd(screen ScreenType) tea.Cmd {
	return func() tea.Msg {
		return screen
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
