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
