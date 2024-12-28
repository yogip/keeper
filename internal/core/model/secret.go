package model

import (
	"encoding/json"
	"fmt"
)

// Generic secret type, secret data stored as JSON marsheled payload
type Secret struct {
	SecretMeta
	Payload []byte
}

func NewSecret(id int64, name string, secretType SecretType, data []byte, note string) *Secret {
	return &Secret{SecretMeta{id, name, secretType, note}, data}
}

// Method to convert generic secret to password object
func (s *Secret) AsPassword() (*Password, error) {
	if s.Type != SecretTypePassword {
		return nil, fmt.Errorf("wrong secret type %s, reqired Password", s.Type)
	}
	var p Password
	err := json.Unmarshal(s.Payload, &p)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Password payload: %w", err)
	}
	p.SecretMeta = s.SecretMeta
	return &p, nil
}

// Method to convert generic secret to password object
func (s *Secret) AsNote() (*Note, error) {
	if s.Type != SecretTypeNote {
		return nil, fmt.Errorf("wrong secret type %s, reqired Note", s.Type)
	}
	var n Note
	err := json.Unmarshal(s.Payload, &n)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Note payload: %w", err)
	}
	n.SecretMeta = s.SecretMeta
	return &n, nil
}

// Method to convert generic secret to password object
func (s *Secret) AsCard() (*Card, error) {
	if s.Type != SecretTypeCard {
		return nil, fmt.Errorf("wrong secret type %s, reqired Card", s.Type)
	}
	var c Card
	err := json.Unmarshal(s.Payload, &c)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Note payload: %w", err)
	}
	c.SecretMeta = s.SecretMeta
	return &c, nil
}

// General secret info
type SecretMeta struct {
	ID   int64
	Name string
	Type SecretType
	Note string `json:"note"`
}

// Password secret
type Password struct {
	SecretMeta
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewPassword(id int64, name, login, password, note string) *Password {
	return &Password{
		SecretMeta: SecretMeta{id, name, SecretTypePassword, note},
		Login:      login,
		Password:   password,
	}
}

// JSON secret date representaion
func (p *Password) GetPayload() ([]byte, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Password payload: %w", err)
	}
	return payload, nil
}

// Text data secret
type Note struct {
	SecretMeta
	Text string `json:"text"`
}

func NewNote(id int64, name, text, note string) *Note {
	return &Note{
		SecretMeta: SecretMeta{id, name, SecretTypeNote, note},
		Text:       text,
	}
}

// JSON secret date representaion
func (p *Note) GetPayload() ([]byte, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Password payload: %w", err)
	}
	return payload, nil
}

// Card secret
type Card struct {
	SecretMeta
	Number     string `json:"number"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	HolderName string `json:"holder_name"`
	CVC        int    `json:"cvc"`
}

func NewCard(id int64, name, number string, month int, year int, holderName string, cvc int, note string) *Card {
	return &Card{
		SecretMeta: SecretMeta{id, name, SecretTypeCard, note},
		Number:     number,
		Month:      month,
		Year:       year,
		HolderName: holderName,
		CVC:        cvc,
	}
}

// JSON secret date representaion
func (c *Card) GetPayload() ([]byte, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Password payload: %w", err)
	}
	return payload, nil
}

func (c *Card) GetDate() string {
	return fmt.Sprintf("%02d/%02d", c.Month, c.Year)
}

// File secret
type File struct {
	SecretMeta
	FileName string `json:"file_name"`
	Body     []byte `json:"file_name"`
}

// JSON secret date representaion
func (c *File) GetPayload() ([]byte, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Password payload: %w", err)
	}
	return payload, nil
}
