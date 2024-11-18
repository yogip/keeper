package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"keeper/internal/core/encryption"
	"keeper/internal/core/model"
	"keeper/internal/infra/repo"
	"keeper/internal/infra/s3"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type SecretService struct {
	repo              *repo.SecretRepo
	encrypter         *encryption.EncryptionService
	s3client          *s3.S3Client
	lastEncKeyVersion int64
}

func NewSecretService(
	repo *repo.SecretRepo,
	s3client *s3.S3Client,
	encrypter *encryption.EncryptionService,
	lastEncKeyVersion int64,
) *SecretService {
	return &SecretService{repo: repo, s3client: s3client, encrypter: encrypter, lastEncKeyVersion: lastEncKeyVersion}
}

func (s *SecretService) ListSecretsMeta(ctx context.Context, req *model.SecretListRequest) (*model.SecretList, error) {
	secrets, err := s.repo.ListSecrets(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get password")
	}
	return &model.SecretList{Secrets: secrets}, nil
}

func (s *SecretService) GetPassword(ctx context.Context, req model.SecretRequest) (*model.Password, error) {
	if req.Type != model.SecretTypePassword {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypePassword, req.Type)
	}

	pwd, err := s.repo.GetPassword(ctx, req)
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

	pwd, err := s.repo.CreatePassword(ctx, req)
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

	pwd, err := s.repo.UpdatePassword(ctx, req)
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

	pwd, err := s.repo.CreateNote(ctx, req)
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

	pwd, err := s.repo.UpdateNote(ctx, req)
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

	note, err := s.repo.GetNote(ctx, req)
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

	encCard, err := s.repo.GetCard(ctx, req)
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

	cardResult, err := s.repo.CreateCard(ctx, encCard, req.UserID)

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

	cardResult, err := s.repo.UpdateCard(ctx, encCard, req.UserID)

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

func (s *SecretService) GetFile(ctx context.Context, req model.SecretRequest) (*model.File, error) {
	if req.Type != model.SecretTypeFile {
		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypeFile, req.Type)
	}

	fileMeta, err := s.repo.GetFileMeta(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get File")
	}

	encText, err := s.s3client.GetObject(ctx, fileMeta.Meta.Path)
	plText, err := s.encrypter.Decrypt(string(encText), fileMeta.DataKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt File")
	}

	return &model.File{Body: plText, FileMeta: *fileMeta.Meta}, nil
}

func (s *SecretService) CreateFile(ctx context.Context, req model.CreateFileRequest) (*model.FileMeta, error) {
	enc, key, err := s.encrypter.Encrypt(req.Body, s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt file")
	}

	buff := bytes.NewBufferString(enc)
	path := uuid.NewString()
	err = s.s3client.PutObject(ctx, path, buff, int64(buff.Len()))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}

	dKey := &model.DataKey{Key: key, Version: s.lastEncKeyVersion}
	meta, err := s.repo.CreateFileMeta(ctx, req, path, dKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}
	return meta.Meta, nil
}
