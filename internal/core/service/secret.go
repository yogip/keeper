package service

import (
	"context"
	"encoding/json"
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

func (s *SecretService) ListSecretsMeta(ctx context.Context, req *model.SecretListRequest) (*model.SecretList, error) {
	secrets, err := s.secretRepo.ListSecrets(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get password")
	}
	return &model.SecretList{Secrets: secrets}, nil
}

func (s *SecretService) GetPassword(ctx context.Context, req model.SecretRequest) (*model.Password, error) {
	if req.Type != model.SecretTypePassword {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypePassword, req.Type)
	}

	pwd, err := s.secretRepo.GetPassword(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get password")
	}

	p, err := s.encrypter.Decrypt(pwd.Item.Password, pwd.DataKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt password")
	}

	pwd.Item.Type = model.SecretTypePassword
	pwd.Item.Password = string(p)
	return pwd.Item, nil
}

func (s *SecretService) GetNote(ctx context.Context, req model.SecretRequest) (*model.Note, error) {
	if req.Type != model.SecretTypeNote {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypeNote, req.Type)
	}

	note, err := s.secretRepo.GetNote(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get note")
	}

	n, err := s.encrypter.Decrypt(note.Item.Note, note.DataKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt note")
	}

	note.Item.Type = model.SecretTypeNote
	note.Item.Note = string(n)
	return note.Item, nil
}

func (s *SecretService) GetCard(ctx context.Context, req model.SecretRequest) (*model.Card, error) {
	if req.Type != model.SecretTypeCard {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypeCard, req.Type)
	}

	encCard, err := s.secretRepo.GetCard(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get card")
	}

	payload, err := s.encrypter.Decrypt(encCard.Payload, encCard.DataKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt card")
	}

	card := model.CardData{}
	err = json.Unmarshal(payload, &card)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal card")
	}

	encCard.Meta.Type = model.SecretTypeCard
	return &model.Card{
		SecretMeta: *encCard.Meta,
		CardData:   card,
	}, nil

}
