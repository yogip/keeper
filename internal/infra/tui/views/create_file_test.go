package views

import (
	"keeper/internal/core/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tea "github.com/charmbracelet/bubbletea"
)

type clientStub struct {
}

func (c *clientStub) CreateSecret(secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error) {
	return &model.Secret{}, nil
}
func (c *clientStub) CreateFileSecret(name, fileName, note string, payload []byte) (int64, error) {
	return 0, nil
}
func (c *clientStub) UpdateSecret(id int64, secretType model.SecretType, name string, note string, payload []byte) (*model.Secret, error) {
	return &model.Secret{}, nil
}
func (c *clientStub) ListSecrets(secretName string) (*model.SecretList, error) {
	return &model.SecretList{}, nil
}
func (c *clientStub) GetSecret(secretID int64) (*model.Secret, error) {
	return &model.Secret{}, nil
}

func TestCreateFileView_Init(t *testing.T) {
	app := clientStub{}
	m := NewCreateFileView(&app)
	cmd := m.Init()

	assert.NotNil(t, cmd)
}

func TestCreateFileView_createSecretCmd(t *testing.T) {
	app := clientStub{}

	f, err := os.CreateTemp("", "tmp_test_file")
	if err != nil {
		require.NoError(t, err)
	}

	defer f.Close()
	defer os.Remove(f.Name())

	m := NewCreateFileView(&app)
	cmd := m.createSecretCmd("test_name", f.Name(), "note")
	msg := cmd()

	expMsg := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	assert.Equal(t, expMsg, msg)
}

func TestCreateFileView_Update(t *testing.T) {
	app := clientStub{}
	m := NewCreateFileView(&app)
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	cmd := m.Update(msg)

	assert.Nil(t, cmd)
	assert.Equal(t, 1, m.focusIndex)
}

func TestCreateFileView_updateInputs(t *testing.T) {
	app := clientStub{}
	m := NewCreateFileView(&app)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	cmd := m.updateInputs(msg)
	assert.NotNil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")}
	cmd = m.updateInputs(msg)
	assert.NotNil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")}
	cmd = m.updateInputs(msg)
	assert.NotNil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}
	cmd = m.updateInputs(msg)
	assert.NotNil(t, cmd)

	assert.Equal(t, "stub", m.nameInput.Value())
}

func TestCreateFileView_View(t *testing.T) {
	expected := "Create Card:\n> Secret name\n\nSelect a file:\n\n  Bummer. No Files Found.\n                         \n                         \n                         \n                         \n\n┃   1 Enter a note                      \n┃                                       \n┃                                       \n┃                                       \n┃                                       \n┃                                       \n\n[ Create ]\n[ Cancel ]\nUse `up` and `down` or `tab` and `shift+tab` to navigate\n"
	app := clientStub{}
	m := NewCreateFileView(&app)

	view := m.View()

	assert.Equal(t, expected, view)
}
