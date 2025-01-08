package views

import (
	"testing"

	"keeper/internal/core/model"
	"keeper/internal/mocks"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSecretView_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockClientApp(ctrl)
	m := NewSecretView(mock)

	secret := model.NewSecret(101, "stub-name", model.SecretTypePassword, []byte("stub-payload"), "stub-note")

	mock.EXPECT().
		GetSecret(gomock.Eq(int64(101))).
		Return(secret, nil).
		AnyTimes()

	cmd := m.Init(secret.ID)
	msg := cmd()
	assert.Equal(t, "", msg)

	assert.Equal(t, secret, m.secret)
}

func TestSecretView_MoveToSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockClientApp(ctrl)
	m := NewSecretView(mock)

	for _, msg := range []tea.KeyMsg{
		{Type: tea.KeyEsc},
		{Type: tea.KeyRunes, Runes: []rune("e")},
	} {
		cmd := m.Update(msg)
		assert.NotNil(t, cmd)

		got := cmd()
		want := changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
		assert.Equal(t, want, got)
	}
}

func TestSecretView_View(t *testing.T) {
	tests := []struct {
		name     string
		secret   *model.Secret
		expected string
	}{
		{
			name:     "Test Password",
			secret:   model.NewSecret(101, "stub-name", model.SecretTypePassword, []byte(`{"login": "stub-login", "password": "stub-password"}`), "stub-note"),
			expected: "Secret ID: 101\nName:      stub-name\n-----------------------------------------------------\nLogin:     stub-login \nPassword:  stub-password \n-----------------------------------------------------\nNote:      stub-note\n\nUse `esc` or `ctr+c` to exit.\nUse `e` to edit secret.",
		},
		{
			name:     "Test Note",
			secret:   model.NewSecret(102, "stub-name", model.SecretTypeNote, []byte(`{"text": "stub-text"}`), "stub-note"),
			expected: "Secret ID: 102\nName:      stub-name\n-----------------------------------------------------\nText:     stub-text \n-----------------------------------------------------\nNote:      stub-note\n\nUse `esc` or `ctr+c` to exit.\nUse `e` to edit secret.",
		},
		{
			name:     "Test Card",
			secret:   model.NewSecret(103, "stub-name", model.SecretTypeCard, []byte(`{"number": "1234 1234 1234 1234", "month": 11, "year": 25, "holder_name": "Holder Name", "cvc": 123}`), "stub-note"),
			expected: "Secret ID: 103\nName:      stub-name\n-----------------------------------------------------\n1234 1234 1234 1234\nHolder Name\nExpired after: 11/25 \t cvc: 123\n-----------------------------------------------------\nNote:      stub-note\n\nUse `esc` or `ctr+c` to exit.\nUse `e` to edit secret.",
		},
		{
			name:     "Test Card",
			secret:   model.NewSecret(103, "stub-name", model.SecretTypeFile, []byte(`{"s3_name":"2_00e63495-1ec9-4b91-ae77-025a535965e6","file_name":"test.txt","file":"ZmlsZSBjb250ZW50"}`), "stub-note"),
			expected: "Secret ID: 103\nName:      stub-name\n-----------------------------------------------------\nFile name:    test.txt\nFile body:\nfile content\n-----------------------------------------------------\nNote:      stub-note\n\nUse `esc` or `ctr+c` to exit.\nUse `e` to edit secret.",
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockClientApp(ctrl)
		m := NewSecretView(mock)
		mock.EXPECT().
			GetSecret(gomock.Eq(tt.secret.ID)).
			Return(tt.secret, nil).
			AnyTimes()

		cmd := m.Init(tt.secret.ID)
		msg := cmd()
		assert.Equal(t, "", msg)

		view := m.View()

		assert.Equal(t, tt.expected, view)
	}
}
