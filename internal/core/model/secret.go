package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Secret struct {
	SecretMeta
	Payload []byte
}

func NewSecret(id int64, name string, secretType SecretType, data []byte) *Secret {
	return &Secret{SecretMeta{id, name, secretType}, data}
}

func (s *Secret) AsPassword() (*Password, error) {
	if s.Type != SecretTypePassword {
		return nil, errors.New("secret is not a password")
	}
	var p Password
	err := json.Unmarshal(s.Payload, &p)
	if err != nil {
		return nil, fmt.Errorf("could not Unmarshal Password payload: %w", err)
	}
	p.SecretMeta = s.SecretMeta
	return &p, nil
}

type SecretMeta struct {
	ID   int64
	Name string
	Type SecretType
	// Note string `json:"note"`
}

type Password struct {
	SecretMeta
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewPassword(id int64, name, login, password, note string) *Password {
	// todo use note
	return &Password{
		SecretMeta: SecretMeta{id, name, SecretTypePassword},
		Login:      login,
		Password:   password,
	}
}

func (p *Password) GetPayload() ([]byte, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("could not Marshal Password payload: %w", err)
	}
	return payload, nil
}
