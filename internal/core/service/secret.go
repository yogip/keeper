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

func (s *SecretService) ListSecretsMeta(ctx context.Context, req model.SecretListRequest) (*model.SecretList, error) {
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

	if secret.Item.Type == model.SecretTypeFile {
		fileMeta, err := secret.Item.AsFile()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get file meta from payload")
		}
		fileBodyEnc, err := s.s3client.GetObject(ctx, fileMeta.S3Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get file from s3")
		}
		fileBody, err := s.encrypter.Decrypt(string(fileBodyEnc), secret.DataKey)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decrypt secret")
		}

		fileMeta.File = fileBody
		secret.Item.Payload, err = fileMeta.GetPayload()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get file payload")
		}
	} else {
		p, err := s.encrypter.Decrypt(string(secret.Item.Payload), secret.DataKey)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decrypt secret")
		}
		secret.Item.Payload = p
	}

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

func (s *SecretService) CreateFile(ctx context.Context, req model.CreateFileRequest) (int64, error) {
	enc, key, err := s.encrypter.Encrypt(req.Payload, s.lastEncKeyVersion)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to encrypt file")
	}

	buff := bytes.NewBufferString(enc)
	s3name := fmt.Sprintf("%d_%s", req.UserID, uuid.NewString())
	err = s.s3client.PutObject(ctx, s3name, buff, int64(buff.Len()))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to save file")
	}

	fileMeta := model.File{
		FileName: req.FileName,
		S3Name:   s3name,
	}
	payload, err := fileMeta.GetPayload()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get payload")
	}

	dKey := &model.DataKey{Key: key, Version: s.lastEncKeyVersion}
	r := model.SecretCreateRequest{
		Name:    req.Name,
		Type:    model.SecretTypeFile,
		Note:    req.Note,
		Payload: payload,
		UserID:  req.UserID,
	}
	sid, err := s.repo.CreateSecret(ctx, &r, dKey)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to crete file")
	}

	return sid, nil
}
