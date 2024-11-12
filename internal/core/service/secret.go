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
	secretRepo        *repo.SecretRepo
	encrypter         *encryption.EncryptionService
	lastEncKeyVersion int64
}

func NewSecretService(repo *repo.SecretRepo, encrypter *encryption.EncryptionService, lastEncKeyVersion int64) *SecretService {
	return &SecretService{secretRepo: repo, encrypter: encrypter, lastEncKeyVersion: lastEncKeyVersion}
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

func (s *SecretService) CreatePassword(ctx context.Context, req model.UpdatePasswordRequest) (*model.Password, error) {
	plaintext := req.Data.Password
	enc, key, err := s.encrypter.Encrypt([]byte(req.Data.Password), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	req.Data.Password = enc
	req.Key = &model.DataKey{Key: key, Version: s.lastEncKeyVersion}

	pwd, err := s.secretRepo.CreatePassword(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}

	pwd.Item.Type = model.SecretTypePassword
	pwd.Item.Password = plaintext
	return pwd.Item, nil
}

func (s *SecretService) UpdatePassword(ctx context.Context, req model.UpdatePasswordRequest) (*model.Password, error) {
	plaintext := req.Data.Password
	enc, key, err := s.encrypter.Encrypt([]byte(req.Data.Password), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	req.Data.Password = enc
	req.Key = &model.DataKey{Key: key, Version: s.lastEncKeyVersion}

	pwd, err := s.secretRepo.UpdatePassword(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update password")
	}

	pwd.Item.Type = model.SecretTypePassword
	pwd.Item.Password = plaintext
	return pwd.Item, nil
}

func (s *SecretService) CreateNote(ctx context.Context, req model.UpdateNoteRequest) (*model.Note, error) {
	plaintext := req.Data.Note
	enc, key, err := s.encrypter.Encrypt([]byte(req.Data.Note), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	req.Data.Note = enc
	req.Key = &model.DataKey{Key: key, Version: s.lastEncKeyVersion}

	pwd, err := s.secretRepo.CreateNote(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}

	pwd.Item.Type = model.SecretTypeNote
	pwd.Item.Note = plaintext
	return pwd.Item, nil
}

func (s *SecretService) UpdateNote(ctx context.Context, req model.UpdateNoteRequest) (*model.Note, error) {
	plaintext := req.Data.Note
	enc, key, err := s.encrypter.Encrypt([]byte(req.Data.Note), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	req.Data.Note = enc
	req.Key = &model.DataKey{Key: key, Version: s.lastEncKeyVersion}

	pwd, err := s.secretRepo.UpdateNote(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update password")
	}

	pwd.Item.Type = model.SecretTypeNote
	pwd.Item.Note = plaintext
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

func (s *SecretService) CreateCard(ctx context.Context, req model.UpdateCardRequest) (*model.Card, error) {
	payload, err := json.Marshal(req.Card.CardData)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete card")
	}

	enc, key, err := s.encrypter.Encrypt([]byte(payload), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	encCard := model.EncryptedCard{
		Meta:    &req.Card.SecretMeta,
		Payload: enc,
		DataKey: &model.DataKey{Key: key, Version: s.lastEncKeyVersion},
	}

	cardResult, err := s.secretRepo.CreateCard(ctx, encCard, req.UserID)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}
	result := model.Card{
		SecretMeta: req.Card.SecretMeta,
		CardData:   req.Card.CardData,
	}
	result.ID = cardResult.Meta.ID
	return &result, nil
}

func (s *SecretService) UpdateCard(ctx context.Context, req model.UpdateCardRequest) (*model.Card, error) {
	payload, err := json.Marshal(req.Card.CardData)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete card")
	}

	enc, key, err := s.encrypter.Encrypt([]byte(payload), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}

	encCard := model.EncryptedCard{
		Meta:    &req.Card.SecretMeta,
		Payload: enc,
		DataKey: &model.DataKey{Key: key, Version: s.lastEncKeyVersion},
	}

	cardResult, err := s.secretRepo.UpdateCard(ctx, encCard, req.UserID)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}
	result := model.Card{
		SecretMeta: req.Card.SecretMeta,
		CardData:   req.Card.CardData,
	}
	result.ID = cardResult.Meta.ID
	return &result, nil
}
