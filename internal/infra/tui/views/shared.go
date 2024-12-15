package views

import (
	"fmt"
	"keeper/internal/core/model"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boldStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Bold(true)
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

type ScreenType string

const (
	ScreenLogin       ScreenType = "login"
	ScreenSignUp      ScreenType = "signup"
	ScreenSecretList  ScreenType = "secret-list"
	ScreenNewSecret   ScreenType = "new-secret"
	ScreenNewPassword ScreenType = "new-secret-pwd"
	ScreenNewNote     ScreenType = "new-secret-note"
	ScreenNewCard     ScreenType = "new-secret-card"
	ScreenNewFile     ScreenType = "new-secret-file"
)

type ClientApp interface {
	ListSecrets(secretName string) (*model.SecretList, error)
}

func changeScreenCmd(screen ScreenType) tea.Cmd {
	log.Println("create changeScreenCmd. screen: ", screen)
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

func blurInput(input *textinput.Model) {
	input.Blur()
	input.PromptStyle = noStyle
	input.TextStyle = noStyle
}

func focusInput(input *textinput.Model) tea.Cmd {
	cmd := input.Focus()
	input.Blur()
	input.PromptStyle = focusedStyle
	input.TextStyle = focusedStyle
	return cmd
}
