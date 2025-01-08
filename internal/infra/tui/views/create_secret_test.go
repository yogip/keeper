package views

import (
	"errors"
	"keeper/internal/core/model"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSecterView_Init(t *testing.T) {
	app := clientStub{}
	m := NewCreateSecretView(&app)
	cmd := m.Init()

	assert.Nil(t, cmd)
}

func TestCreateSecterView_nextStageView(t *testing.T) {
	app := clientStub{}
	m := NewCreateSecretView(&app)

	tests := []struct {
		name           string
		argsSecretType model.SecretType
		want           tea.Msg
	}{
		{
			name:           "Test password",
			argsSecretType: model.SecretTypePassword,
			want:           changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertPassword})(),
		},
		{
			name:           "Test note",
			argsSecretType: model.SecretTypeNote,
			want:           changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertNote})(),
		},
		{
			name:           "Test card",
			argsSecretType: model.SecretTypeCard,
			want:           changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertCard})(),
		},
		{
			name:           "Test file",
			argsSecretType: model.SecretTypeFile,
			want:           changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertFile})(),
		},
		{
			name:           "Test unknown",
			argsSecretType: model.SecretType("unknown"),
			want:           ErrorCmd(errors.New("Select secret type"), time.Second*15)(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.secretType = tt.argsSecretType
			gotCmd := m.nextStageView()
			msg := gotCmd()
			if e, ok := msg.(ErrorMsg); ok {
				wantE, ok := tt.want.(ErrorMsg)
				require.True(t, ok)
				assert.Equal(t, e.Err, wantE.Err)
			} else {
				assert.Equal(t, tt.want, msg)
			}

		})
	}
}

func TestCreateSecterView_Update(t *testing.T) {
	app := clientStub{}
	m := NewCreateSecretView(&app)
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	cmd := m.Update(msg)
	assert.NotNil(t, cmd)

	got := cmd()
	want := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenUpsertPassword})()

	assert.Equal(t, want, got)
}

func TestCreateSecterView_UpdateDown(t *testing.T) {
	app := clientStub{}
	m := NewCreateSecretView(&app)
	msg := tea.KeyMsg{Type: tea.KeyDown}

	// go down thourought the all secret types
	for i := 0; i < 3; i++ {
		cmd := m.Update(msg)
		assert.Nil(t, cmd)
		assert.Equal(t, 0, m.focusIndex)
	}

	// select next button
	cmd := m.Update(msg)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, m.focusIndex)

	// select cancel button
	cmd = m.Update(msg)
	assert.Nil(t, cmd)
	assert.Equal(t, 2, m.focusIndex)
}

func TestCreateSecterView_View(t *testing.T) {
	expected := "  Select secret type?  \n                       \n  > password           \n    note               \n    card               \n    file               \n                       \nSecret type: \n             \n\n[ Next ]\n[ Cancel ]\nUse `up` and `down` or `tab` and `shift+tab` to navigate\n"
	app := clientStub{}
	m := NewCreateSecretView(&app)

	view := m.View()

	assert.Equal(t, expected, view)
}
