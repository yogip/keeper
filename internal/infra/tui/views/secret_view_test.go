package views

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestSecretView_loginCmd(t *testing.T) {
// 	app := clientStub{}
// 	m := NewSecretView(&app)

// 	cmd := m.loginCmd("user", "password")
// 	msg := cmd()

// 	expMsg := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
// 	assert.Equal(t, expMsg, msg)
// }

// func TestSecretView_MoveToSignUp(t *testing.T) {
// 	app := clientStub{}
// 	m := NewSecretView(&app)
// 	msg := tea.KeyMsg{Type: tea.KeyEnter}

// 	m.focusIndex = m.focusSignUp
// 	cmd := m.Update(msg)
// 	assert.NotNil(t, cmd)

// 	got := cmd()
// 	want := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSignUp})()

// 	assert.Equal(t, want, got)
// }

// func TestSecretView_UpdateDown(t *testing.T) {
// 	app := clientStub{}
// 	m := NewSecretView(&app)
// 	msg := tea.KeyMsg{Type: tea.KeyDown}

// 	// go down to password input
// 	cmd := m.Update(msg)
// 	assert.NotNil(t, cmd)
// 	assert.Equal(t, 1, m.focusIndex)

// 	// go down to submit button
// 	cmd = m.Update(msg)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 2, m.focusIndex)

// 	// go down to sign up button
// 	cmd = m.Update(msg)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, 3, m.focusIndex)
// }

// func TestSecretView_View(t *testing.T) {
// 	expected := "Enter your credentials:\n\n> Email\n> Password\n\n[ Submit ]\n[ Sign Up ]\nUse `up` and `down` or `tab` and `shift+tab` to navigate."
// 	app := clientStub{}
// 	m := NewSecretView(&app)

// 	view := m.View()

// 	assert.Equal(t, expected, view)
// }
