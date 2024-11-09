package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	coreErrors "keeper/internal/core/errors"
	"keeper/internal/core/model"
	"keeper/internal/logger"
	"keeper/internal/retrier"
)

type SecretRepo struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

func NewSecretRepo(db *sql.DB) *SecretRepo {
	ret := &retrier.Retrier{
		Strategy: retrier.Backoff(
			3,             // max attempts
			1*time.Second, // initial delay
			3,             // multiplier
			5*time.Second, // max delay
		),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Log.Debug(fmt.Sprintf("Retrying DB. retry #%d: %v", n, err))
		},
	}

	repo := &SecretRepo{db: db, retrier: ret}

	logger.Log.Info("SecretRepo initialized")
	return repo
}

func (r *SecretRepo) ListSecrets(ctx context.Context, user *model.User) (*model.User, error) {
	return &model.User{}, nil
}

// Read Password with secret key from DB.
func (r *SecretRepo) GetPassword(ctx context.Context, req model.SecretRequest) (*model.Password, *model.DataKey, error) {
	pwd := &model.Password{}
	dataSecret := &model.DataKey{}
	query := `
		SELECT 
			p.id, p.name, t.folder_id, f.name as folder, login, password, sc_version, sc
		FROM 
			public.passwords t
		LEFT JOIN 
			folders f on f.id = t.folder_id AND f.user_id = t.user_id
		WHERE t.id = $1 AND t.user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)
		err := row.Scan(
			&pwd.ID,
			&pwd.Name,
			&pwd.Folder.ID,
			&pwd.Folder.Name,
			&pwd.Login,
			&pwd.Password,
			&dataSecret.Version,
			&dataSecret.Key,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading password: %w", err)
	}
	return pwd, dataSecret, nil
}

// Read Note with secret key from DB.
func (r *SecretRepo) GetNote(ctx context.Context, req model.SecretRequest) (*model.Note, *model.DataKey, error) {
	note := &model.Note{}
	dataSecret := &model.DataKey{}
	query := `
		SELECT 
			t.id, t.name, t.folder_id, f.name as folder, note, sc_version, sc
		FROM 
			public.notes t
		LEFT JOIN 
			folders f on f.id = p.folder_id AND f.user_id = t.user_id
		WHERE t.id = $1 AND t.user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)
		err := row.Scan(
			&note.ID,
			&note.Name,
			&note.Folder.ID,
			&note.Folder.Name,
			&note.Note,
			&dataSecret.Version,
			&dataSecret.Key,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading password: %w", err)
	}
	return note, dataSecret, nil
}
