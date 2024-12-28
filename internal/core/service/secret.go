package service

import (
	"bytes"
	"context"
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
		return nil, errors.Wrapf(err, "failed to list secrets")
	}
	return &model.SecretList{Secrets: secrets}, nil
}

func (s *SecretService) GetSecret(ctx context.Context, req model.SecretRequest) (*model.Secret, error) {
	secret, err := s.repo.GetSecret(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get secret")
	}

	p, err := s.encrypter.Decrypt(string(secret.Item.Payload), secret.DataKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt secret")
	}

	secret.Item.Payload = p
	return secret.Item, nil
}

func (s *SecretService) CreateSecret(ctx context.Context, req model.SecretCreateRequest) (*model.Secret, error) {
	resp := model.Secret{
		SecretMeta: model.SecretMeta{Name: req.Name, Type: req.Type, Note: req.Note},
		Payload:    req.Payload,
	}
	enc, key, err := s.encrypter.Encrypt(req.Payload, s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt secret")
	}

	req.Payload = []byte(enc)
	secretID, err := s.repo.CreateSecret(ctx, &req, &model.DataKey{Key: key, Version: s.lastEncKeyVersion})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create secret")
	}

	resp.ID = secretID
	return &resp, nil
}

func (s *SecretService) UpdateSecret(ctx context.Context, req model.SecretUpdateRequest) (*model.Secret, error) {
	resp := model.Secret{
		SecretMeta: model.SecretMeta{ID: req.ID, Name: req.Name, Type: req.Type, Note: req.Note},
		Payload:    req.Payload,
	}
	enc, key, err := s.encrypter.Encrypt([]byte(req.Payload), s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt password")
	}
	req.Payload = []byte(enc)

	err = s.repo.UpdateSecret(ctx, &req, &model.DataKey{Key: key, Version: s.lastEncKeyVersion})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update password")
	}

	return &resp, nil
}

// todo ????
// func (s *SecretService) GetFile(ctx context.Context, req model.SecretRequest) (*model.File, error) {
// 	if req.Type != model.SecretTypeFile {
// 		return nil, fmt.Errorf("type must be %s, got: %s", model.SecretTypeFile, req.Type)
// 	}

// 	fileMeta, err := s.repo.GetFileMeta(ctx, req)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to get File")
// 	}

// 	encText, err := s.s3client.GetObject(ctx, fileMeta.Meta.Path)
// 	plText, err := s.encrypter.Decrypt(string(encText), fileMeta.DataKey)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "failed to decrypt File")
// 	}

// 	return &model.File{Body: plText, FileMeta: *fileMeta.Meta}, nil
// }

func (s *SecretService) CreateFile(ctx context.Context, req model.CreateFileRequest) (*model.FileMeta, error) {
	enc, key, err := s.encrypter.Encrypt(req.Payload, s.lastEncKeyVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt file")
	}

	buff := bytes.NewBufferString(enc)
	s3name := fmt.Sprintf("%d_%s", req.UserID, uuid.NewString())
	err = s.s3client.PutObject(ctx, s3name, buff, int64(buff.Len()))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save file")
	}

	dKey := &model.DataKey{Key: key, Version: s.lastEncKeyVersion}
	r := model.SecretCreateRequest{
		Name:    req.Name,
		Type:    model.SecretTypeFile,
		Note:    req.Note,
		Payload: []byte(enc), // todo !!!
		UserID:  req.UserID,
	}
	sid, err := s.repo.CreateSecret(ctx, &r, dKey)
	// meta, err := s.repo.CreateSecret(ctx, req, s3name, dKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to crete password")
	}

	resp := model.FileMeta{
		SecretMeta: model.SecretMeta{
			ID: sid,
		},
	}
	return &resp, nil
}
