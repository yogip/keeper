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

type EncryptedPassword struct {
	Item    *model.Password
	DataKey *model.DataKey
}

type EncryptedNote struct {
	Item    *model.Note
	DataKey *model.DataKey
}

type EncryptedCard struct {
	Payload string
	Meta    *model.SecretMeta
	DataKey *model.DataKey
}

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

func (r *SecretRepo) ListSecrets(ctx context.Context, req *model.SecretListRequest) ([]*model.SecretMeta, error) {
	items := make([]*model.SecretMeta, 0, 10)
	var query string
	switch req.Type {
	case model.SecretTypePassword:
		query = "SELECT id, name FROM public.passwords WHERE user_id = $1"
	case model.SecretTypeNote:
		query = "SELECT id, name FROM public.notes WHERE user_id = $1"
	case model.SecretTypeCard:
		query = "SELECT id, name FROM public.cards WHERE user_id = $1"
	}

	fun := func() error {
		rows, err := r.db.QueryContext(ctx, query, req.UserID)
		if err != nil {
			return fmt.Errorf("selecting ListSecrets error: %w", err)
		}

		for rows.Next() {
			m := model.SecretMeta{Type: req.Type}

			err = rows.Scan(&m.ID, &m.Name)
			if err != nil {
				return fmt.Errorf("read ListSecrets error: %w", err)
			}

			items = append(items, &m)
		}

		err = rows.Err()
		if err != nil {
			return fmt.Errorf("ListSecrets error: %w", err)
		}
		return nil
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Read Password with secret key from DB.
func (r *SecretRepo) GetPassword(ctx context.Context, req model.SecretRequest) (*EncryptedPassword, error) {
	pwd := &model.Password{}
	dataSecret := &model.DataKey{}
	query := `
		SELECT 
			p.id, p.name, login, password, sc_version, sc
		FROM 
			public.passwords p
		WHERE p.id = $1 AND p.user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)

		err := row.Scan(
			&pwd.ID,
			&pwd.Name,
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
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	return &EncryptedPassword{Item: pwd, DataKey: dataSecret}, nil
}

// Read Note with secret key from DB.
func (r *SecretRepo) GetNote(ctx context.Context, req model.SecretRequest) (*EncryptedNote, error) {
	note := &model.Note{}
	dataSecret := &model.DataKey{}
	query := `
		SELECT 
			n.id, n.name, note, sc_version, sc
		FROM 
			public.notes n
		WHERE n.id = $1 AND n.user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)

		err := row.Scan(
			&note.ID,
			&note.Name,
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
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	return &EncryptedNote{Item: note, DataKey: dataSecret}, nil
}

// Read Card encrypted data with secret key from DB.
func (r *SecretRepo) GetCard(ctx context.Context, req model.SecretRequest) (*EncryptedCard, error) {
	card := EncryptedCard{
		Meta:    &model.SecretMeta{},
		DataKey: &model.DataKey{},
	}
	query := `
		SELECT 
			n.id, n.name, payload, sc_version, sc
		FROM 
			public.cards n
		WHERE n.id = $1 AND n.user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)

		err := row.Scan(
			&card.Meta.ID,
			&card.Meta.Name,
			&card.Payload,
			&card.DataKey.Version,
			&card.DataKey.Key,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	return &card, nil
}
