package views

import (
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SecretView struct {
	app    ClientApp
	secret *model.Secret
}

func NewSecretView(app ClientApp) *SecretView {
	return &SecretView{
		app:    app,
		secret: nil,
	}
}

func (m *SecretView) Init(secretID int64) tea.Cmd {
	return func() tea.Msg {
		secret, err := m.app.GetSecret(secretID)
		if err != nil {
			log.Println("Error loading secret", err)
			return NewErrorMsg(err, time.Second*30)
		}

		m.secret = secret
		return ""
	}
}

func (m *SecretView) Update(msg tea.Msg) tea.Cmd {
	key, ok := msg.(tea.KeyMsg)
	if ok && (key.String() == "ctrl+c" || key.String() == "esc") {
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
	}
	if ok && key.String() == "e" {
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
	}
	return nil
}

func (m *SecretView) View() string {
	if m.secret == nil {
		return "Loading secret..."
	}

	var b strings.Builder
	b.WriteString(boldStyle.Render("Secret ID:") + fmt.Sprintf(" %d", m.secret.ID))
	b.WriteRune('\n')

	b.WriteString(boldStyle.Render("Name: \t") + fmt.Sprintf(" %s", m.secret.Name))
	b.WriteRune('\n')
	b.WriteString("-----------------------------------------------------\n")

	// Secret Body
	switch m.secret.Type {
	case model.SecretTypePassword:
		p, err := m.secret.AsPassword()
		if err != nil {
			b.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		b.WriteString(boldStyle.Render("Login:\t") + fmt.Sprintf(" %s \n", p.Login))
		b.WriteString(boldStyle.Render("Password: ") + fmt.Sprintf(" %s \n", p.Password))
	case model.SecretTypeNote:
		n, err := m.secret.AsNote()
		if err != nil {
			b.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		b.WriteString(boldStyle.Render("Text:\t") + fmt.Sprintf(" %s \n", n.Text))
	case model.SecretTypeCard:
		c, err := m.secret.AsCard()
		if err != nil {
			b.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		b.WriteString(c.Number)
		b.WriteRune('\n')
		b.WriteString(c.HolderName)
		b.WriteRune('\n')
		b.WriteString(fmt.Sprintf("Expired after: %s \t cvc: %d", c.GetDate(), c.CVC))
		b.WriteRune('\n')

	}

	b.WriteString("-----------------------------------------------------\n")
	b.WriteString(boldStyle.Render("Note: \t") + fmt.Sprintf(" %s", m.secret.Note))
	b.WriteRune('\n')

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `esc` or `ctr+c` to exit."))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `e` to edit secret."))

	return b.String()
}
