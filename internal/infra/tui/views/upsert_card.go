package views

import (
	"errors"
	"fmt"
	"keeper/internal/core/model"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type UpsertCardView struct {
	focusIndex  int
	focusMax    int
	focusName   int
	focusNumber int
	focusHolder int
	focusDate   int
	focusCVC    int
	focusNote   int
	focusSubmit int
	focusCancel int

	nameInput   textinput.Model
	numberInput textinput.Model
	holderInput textinput.Model
	dateInput   textinput.Model
	cvcInput    textinput.Model
	noteInput   textarea.Model

	secretID *int64
	app      ClientApp
}

func cardNumberValidator(s string) error {
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("Card number is too long")
	}

	if len(s) < 16+3 {
		return fmt.Errorf("Card number is too sort")
	}

	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	// The 3 character should be a slash (/)
	// The rest should be numbers
	if len(s) != 5 || strings.Index(s, "/") != 2 {
		return fmt.Errorf("Expiration date should has format MM/YY")
	}

	splited := strings.Split(s, "/")
	for i, v := range splited {
		part := "Month"
		if i == 1 {
			part = "Year"
		}
		_, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("Expiration date is invalid - %s has to be a number", part)
		}
	}

	return nil
}

func cvcValidator(s string) error {
	// The CVV should be a number of 3 digits
	if len(s) != 3 {
		return fmt.Errorf("cvc is too short")
	}

	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("cvc has to be a number")
	}
	return nil
}

func NewUpsertCardView(app ClientApp) *UpsertCardView {
	// Name
	nameInput := textinput.New()
	nameInput.CharLimit = 50
	nameInput.Placeholder = "Secret name"
	nameInput.Focus()
	nameInput.Cursor.Style = cursorStyle
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	// Card Data
	numberInput := textinput.New()
	numberInput.Placeholder = "1234 **** **** 1234"
	numberInput.Blur()
	numberInput.CharLimit = 19
	numberInput.Width = 30
	numberInput.Cursor.Style = cursorStyle
	numberInput.PromptStyle = blurredStyle
	numberInput.TextStyle = blurredStyle
	numberInput.Validate = cardNumberValidator

	holderInput := textinput.New()
	holderInput.Placeholder = "Holder Name"
	holderInput.CharLimit = 20
	holderInput.Width = 20
	holderInput.PromptStyle = blurredStyle
	holderInput.Cursor.Style = cursorStyle
	holderInput.TextStyle = blurredStyle

	dateInput := textinput.New()
	dateInput.Placeholder = "MM/YY"
	dateInput.CharLimit = 5
	dateInput.Width = 5
	dateInput.Prompt = "date: "
	dateInput.PromptStyle = blurredStyle
	dateInput.Cursor.Style = cursorStyle
	dateInput.TextStyle = blurredStyle
	dateInput.Validate = expValidator

	cvcInput := textinput.New()
	cvcInput.Placeholder = "CVC"
	cvcInput.CharLimit = 3
	cvcInput.Width = 3
	cvcInput.Prompt = "cvc: "
	cvcInput.PromptStyle = blurredStyle
	cvcInput.Cursor.Style = cursorStyle
	cvcInput.TextStyle = blurredStyle
	cvcInput.Validate = cvcValidator

	// Note
	noteInput := textarea.New()
	noteInput.Cursor.Style = cursorStyle
	noteInput.CharLimit = 50
	noteInput.Placeholder = "Enter a note"
	noteInput.Blur()

	return &UpsertCardView{
		focusIndex: 0,
		focusMax:   8,

		focusName:   0,
		focusNumber: 1,
		focusHolder: 2,
		focusDate:   3,
		focusCVC:    4,
		focusNote:   5,
		focusSubmit: 6,
		focusCancel: 7,

		nameInput:   nameInput,
		numberInput: numberInput,
		holderInput: holderInput,
		dateInput:   dateInput,
		cvcInput:    cvcInput,
		noteInput:   noteInput,

		app: app,
	}
}

func (m *UpsertCardView) Init(secretID *int64) tea.Cmd {
	return func() tea.Msg {
		if secretID == nil {
			return ""
		}
		log.Println("Init secret edit view:", *secretID)
		secret, err := m.app.GetSecret(*secretID)
		if err != nil {
			log.Println("Error loading secret", err)
			return NewErrorMsg(err, time.Second*30)
		}
		card, err := secret.AsCard()
		if err != nil {
			log.Println("Could not open secret for editing", err)
			return NewErrorMsg(err, time.Second*30)
		}

		m.secretID = &card.ID
		m.nameInput.SetValue(card.Name)
		m.noteInput.SetValue(card.Note)

		m.numberInput.SetValue(card.Number)
		m.holderInput.SetValue(card.HolderName)
		m.dateInput.SetValue(card.GetDate())
		m.cvcInput.SetValue(fmt.Sprint(card.CVC))
		log.Println("Secret data loaded:", card.ID, card.Name)
		return ""
	}
}

func (m *UpsertCardView) upsertSecretCmd(name, number, holderName string, month, year int, cvc int, note string) tea.Cmd {
	return func() tea.Msg {
		s := model.NewCard(0, name, number, month, year, holderName, cvc, note)
		payload, err := s.GetPayload()
		if err != nil {
			log.Println("creating note payload error", err)
			return NewErrorMsg(err, time.Second*10)
		}

		if m.secretID == nil {
			_, err = m.app.CreateSecret(model.SecretTypeCard, name, note, payload)
		} else {
			_, err = m.app.UpdateSecret(*m.secretID, model.SecretTypeCard, name, note, payload)
		}
		if err != nil {
			log.Println("call grpc method error", err)
			return NewErrorMsg(err, time.Second*10)
		}
		log.Println("Succesfully upsert secret", name)
		return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})()
	}
}

func (m *UpsertCardView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenNewSecret})
		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// Allow to do a line break for text areas elements
			if s == "enter" && m.focusIndex == m.focusNote {
				return m.updateInputs(msg)
			}

			if s == "enter" && m.focusIndex == m.focusCancel {
				return changeScreenCmd(&ScreenTypeMsg{Screen: ScreenSecretList})
			}
			if s == "enter" && m.focusIndex == m.focusSubmit {
				cmds := make([]tea.Cmd, 0, 5)
				if m.nameInput.Value() == "" {
					cmds = append(cmds, ErrorCmd(errors.New("Secret Name cannot be empty"), time.Second*15))
				}
				if m.numberInput.Err != nil {
					errCmd := ErrorCmd(fmt.Errorf("Card number validation error: %w", m.numberInput.Err), time.Second*15)
					cmds = append(cmds, errCmd)
				}
				if m.holderInput.Err != nil {
					errCmd := ErrorCmd(fmt.Errorf("Holder Name validation error: %w", m.holderInput.Err), time.Second*15)
					cmds = append(cmds, errCmd)
				}
				if m.dateInput.Err != nil {
					errCmd := ErrorCmd(fmt.Errorf("Card expiration date validation error: %w", m.dateInput.Err), time.Second*15)
					cmds = append(cmds, errCmd)
				}
				if m.cvcInput.Err != nil {
					errCmd := ErrorCmd(fmt.Errorf("Card CVC validation error: %w", m.cvcInput.Err), time.Second*15)
					cmds = append(cmds, errCmd)
				}
				if len(cmds) > 0 {
					return tea.Batch(cmds...)
				}

				// Following errors never must be returned
				di := strings.Split(m.dateInput.Value(), "/")
				// mm, err := strconv.ParseInt(di[0], 10, 64)
				month, err := strconv.Atoi(di[0])
				if err != nil {
					return ErrorCmd(fmt.Errorf("Card expiration date error: %w", err), time.Second*10)
				}
				year, err := strconv.Atoi(di[1])
				if err != nil {
					return ErrorCmd(fmt.Errorf("Card expiration date error: %w", err), time.Second*10)
				}
				cvc, err := strconv.Atoi(m.cvcInput.Value())
				if err != nil {
					return ErrorCmd(fmt.Errorf("Card CVC error: %w", err), time.Second*10)
				}
				return m.upsertSecretCmd(
					m.nameInput.Value(),
					m.numberInput.Value(),
					m.holderInput.Value(),
					month,
					year,
					cvc,
					m.noteInput.Value(),
				)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.focusMax {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.focusMax
			}

			var cmd tea.Cmd
			// Name Component
			if m.focusIndex == m.focusName {
				cmd = m.nameInput.Focus()
				m.nameInput.PromptStyle = focusedStyle
				m.nameInput.TextStyle = focusedStyle
			} else {
				m.nameInput.Blur()
				m.nameInput.PromptStyle = noStyle
				m.nameInput.TextStyle = noStyle
			}

			// Card Number Component
			if m.focusIndex == m.focusNumber {
				cmd = m.numberInput.Focus()
				m.numberInput.PromptStyle = focusedStyle
				m.numberInput.TextStyle = focusedStyle
			} else {
				m.numberInput.Blur()
				m.numberInput.PromptStyle = noStyle
				m.numberInput.TextStyle = noStyle
			}

			// Holder Name Component
			if m.focusIndex == m.focusHolder {
				cmd = m.holderInput.Focus()
				m.holderInput.PromptStyle = focusedStyle
				m.holderInput.TextStyle = focusedStyle
			} else {
				m.holderInput.Blur()
				m.holderInput.PromptStyle = noStyle
				m.holderInput.TextStyle = noStyle
			}

			// Date Component
			if m.focusIndex == m.focusDate {
				cmd = m.dateInput.Focus()
				m.dateInput.PromptStyle = focusedStyle
				m.dateInput.TextStyle = focusedStyle
			} else {
				m.dateInput.Blur()
				m.dateInput.PromptStyle = noStyle
				m.dateInput.TextStyle = noStyle
			}

			// CVC Component
			if m.focusIndex == m.focusCVC {
				cmd = m.cvcInput.Focus()
				m.cvcInput.PromptStyle = focusedStyle
				m.cvcInput.TextStyle = focusedStyle
			} else {
				m.cvcInput.Blur()
				m.cvcInput.PromptStyle = noStyle
				m.cvcInput.TextStyle = noStyle
			}

			// Note Component
			if m.focusIndex == m.focusNote {
				cmd = m.noteInput.Focus()
			} else {
				m.noteInput.Blur()
			}

			return cmd
		}

		return m.updateInputs(msg)
	}
	return nil
}

func (m *UpsertCardView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 6)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.noteInput, cmds[1] = m.noteInput.Update(msg)

	m.numberInput, cmds[2] = m.numberInput.Update(msg)
	m.holderInput, cmds[3] = m.holderInput.Update(msg)
	m.dateInput, cmds[4] = m.dateInput.Update(msg)
	m.cvcInput, cmds[5] = m.cvcInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m *UpsertCardView) View() string {
	var b strings.Builder
	action := "Create"
	if m.secretID != nil {
		action = "Update"
	}
	b.WriteString(boldStyle.Render(action + " Card:"))
	b.WriteRune('\n')

	b.WriteString(m.nameInput.View())
	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(regularStyle.Render("Card data:"))
	b.WriteRune('\n')
	b.WriteString(m.numberInput.View())
	b.WriteRune('\n')

	b.WriteString(m.holderInput.View())
	b.WriteRune('\n')

	b.WriteString(m.dateInput.View())
	b.WriteString("         ")
	b.WriteString(m.cvcInput.View())
	b.WriteRune('\n')

	b.WriteRune('\n')
	b.WriteString(m.noteInput.View())
	b.WriteRune('\n')

	// submit button
	b.WriteString("\n")
	button := blurredStyle.Render(fmt.Sprintf("[ %s ]", action))
	if m.focusIndex == m.focusSubmit {
		button = fmt.Sprintf("[ %s ]", focusedStyle.Render(action))
	}
	fmt.Fprintf(&b, "%s", button)

	// cancel button
	b.WriteString("\n")
	cancelBtn := blurredStyle.Render("[ Cancel ]")
	if m.focusIndex == m.focusCancel {
		cancelBtn = fmt.Sprintf("[ %s ]", focusedStyle.Render("Cancel"))
	}
	fmt.Fprintf(&b, "%s", cancelBtn)

	// help info
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Use `up` and `down` or `tab` and `shift+tab` to navigate"))
	b.WriteString("\n")

	return b.String()
}
