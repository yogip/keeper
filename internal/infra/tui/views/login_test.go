package views

import (
	"keeper/internal/core/model"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

type iamStub struct {
}

func (c *iamStub) Login(req model.UserRequest) error  { return nil }
func (c *iamStub) SignUp(req model.UserRequest) error { return nil }

func TestLoginView_loginCmd(t *testing.T) {
	iam := iamStub{}
	m := NewLoginView(&iam)

	cmd := m.loginCmd("user", "password")
	msg := cmd()

	expMsg := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	assert.Equal(t, expMsg, msg)
}

func TestLoginView_MoveToSignUp(t *testing.T) {
	iam := iamStub{}
	m := NewLoginView(&iam)
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	m.focusIndex = m.focusSignUp
	cmd := m.Update(msg)
	assert.NotNil(t, cmd)

	got := cmd()
	want := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSignUp})()

	assert.Equal(t, want, got)
}

func TestLoginView_UpdateDown(t *testing.T) {
	iam := iamStub{}
	m := NewLoginView(&iam)
	msg := tea.KeyMsg{Type: tea.KeyDown}

	// go down to password input
	cmd := m.Update(msg)
	assert.NotNil(t, cmd)
	assert.Equal(t, 1, m.focusIndex)

	// go down to submit button
	cmd = m.Update(msg)
	assert.Nil(t, cmd)
	assert.Equal(t, 2, m.focusIndex)

	// go down to sign up button
	cmd = m.Update(msg)
	assert.Nil(t, cmd)
	assert.Equal(t, 3, m.focusIndex)
}

func TestLoginView_View(t *testing.T) {
	expected := "Enter your credentials:\n\n> Email\n> Password\n\n[ Submit ]\n[ Sign Up ]\nUse `up` and `down` or `tab` and `shift+tab` to navigate."
	iam := iamStub{}
	m := NewLoginView(&iam)

	view := m.View()

	assert.Equal(t, expected, view)
}
