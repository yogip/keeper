package service

import (
	"context"
	"fmt"
	"keeper/internal/core/encryption"
	"keeper/internal/core/model"
	"keeper/internal/infra/repo"

	"github.com/go-faster/errors"
)

type SecretService struct {
	secretRepo *repo.SecretRepo
	encrypter  *encryption.EncryptionService
}

func NewSecretService(repo *repo.SecretRepo, encrypter *encryption.EncryptionService) *SecretService {
	return &SecretService{secretRepo: repo, encrypter: encrypter}
}

func (s *SecretService) LisSecretsMeta(*model.SecretListRequest) (*model.SecretList, error) {
	return &model.SecretList{}, nil
}

func (s *SecretService) GetPassword(ctx context.Context, req model.SecretRequest) (*model.Password, error) {
	if req.Type != model.SecretTypePassword {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypePassword, req.Type)
	}

	pwd, key, err := s.secretRepo.GetPassword(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get password")
	}

	p, err := s.encrypter.Decrypt([]byte(pwd.Password), key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt password")
	}

	pwd.Password = string(p)
	return pwd, nil
}

func (s *SecretService) GetNote(ctx context.Context, req model.SecretRequest) (*model.Note, error) {
	if req.Type != model.SecretTypeNote {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypeNote, req.Type)
	}

	note, key, err := s.secretRepo.GetNote(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get password")
	}

	n, err := s.encrypter.Decrypt([]byte(note.Note), key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt password")
	}

	note.Note = string(n)
	return note, nil
}
