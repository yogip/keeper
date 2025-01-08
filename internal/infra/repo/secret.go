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

type EncryptedSecret struct {
	Item    *model.Secret
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

func (r *SecretRepo) ListSecrets(ctx context.Context, req model.SecretListRequest) ([]*model.SecretMeta, error) {
	items := make([]*model.SecretMeta, 0, 10)
	query := `
		SELECT id, name, secret_type, note 
		FROM public.secrets WHERE user_id = $1 AND name ^@ $2 
		ORDER BY secret_type, name;
	`

	fun := func() error {
		rows, err := r.db.QueryContext(ctx, query, req.UserID, req.Name)
		if err != nil {
			return fmt.Errorf("selecting ListSecrets error: %w", err)
		}

		for rows.Next() {
			var m model.SecretMeta

			err = rows.Scan(&m.ID, &m.Name, &m.Type, &m.Note)
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

// Read Secret with secret key from DB.
func (r *SecretRepo) GetSecret(ctx context.Context, req model.SecretRequest) (*EncryptedSecret, error) {
	secret := &model.Secret{}
	dataSecret := &model.DataKey{}
	query := `
		SELECT 
			id, name, payload, note, secret_type, sc_version, sc
		FROM 
			secrets
		WHERE id = $1 AND user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)

		err := row.Scan(
			&secret.ID,
			&secret.Name,
			&secret.Payload,
			&secret.Note,
			&secret.Type,
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
		return nil, fmt.Errorf("error reading secret: %w", err)
	}
	return &EncryptedSecret{Item: secret, DataKey: dataSecret}, nil
}

func (r *SecretRepo) CreateSecret(ctx context.Context, req *model.SecretCreateRequest, key *model.DataKey) (int64, error) {
	var secretID int64
	query := `
	INSERT INTO 
		secrets(user_id, name, note, payload, secret_type, sc_version, sc) 
	values($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			req.UserID, req.Name, req.Note, req.Payload, req.Type,
			key.Version,
			key.Key,
		)
		err := row.Scan(&secretID)
		if err != nil {
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return 0, fmt.Errorf("create secret error: %w", err)
	}
	return secretID, nil
}

func (r *SecretRepo) UpdateSecret(ctx context.Context, req *model.SecretUpdateRequest, key *model.DataKey) error {
	query := "UPDATE secrets SET name=$1, note=$2, payload=$3, sc_version=$4, sc=$5 WHERE id=$6 AND user_id=$7;"

	fun := func() error {
		result, err := r.db.ExecContext(
			ctx, query,
			req.Name, req.Note, req.Payload,
			key.Version, key.Key,
			req.ID, req.UserID,
		)
		if err != nil {
			return err
		}
		affeted, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affeted == 0 {
			return coreErrors.ErrNotFound404
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("update secret error: %w", err)
	}
	return nil
}
