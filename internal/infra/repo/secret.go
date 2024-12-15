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

type EncryptedFileMeta struct {
	Meta    *model.FileMeta
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
	query := `
		SELECT id, name, secret_type 
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

			err = rows.Scan(&m.ID, &m.Name, &m.Type)
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
			id, name, payload, secret_type, sc_version, sc
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
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	return &EncryptedSecret{Item: secret, DataKey: dataSecret}, nil
}

func (r *SecretRepo) CreateSecret(ctx context.Context, req *model.SecretCreateRequest, key *model.DataKey) (int64, error) {
	var secretID int64
	query := `
	INSERT INTO 
		secrets(user_id, name, payload, secret_type, sc_version, sc) 
	values($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			req.UserID, req.Name, req.Payload, req.Type,
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
		return 0, fmt.Errorf("create password error: %w", err)
	}
	return secretID, nil
}

func (r *SecretRepo) UpdateSecret(ctx context.Context, req *model.SecretUpdateRequest, key *model.DataKey) error {
	query := "UPDATE secrets SET name=$1, type=$2, payload=$3, sc_version=$4, sc=$5 WHERE id=$6 AND user_id=$7;"

	fun := func() error {
		result, err := r.db.ExecContext(
			ctx, query,
			req.Name, req.Type, req.Payload,
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

func (r *SecretRepo) CreatePassword(ctx context.Context, req model.UpdatePasswordRequest) (*EncryptedPassword, error) {
	result := req.Data
	query := `
	INSERT INTO 
		passwords(user_id, name, login, password, sc_version, sc) 
	values($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			req.UserID, req.Data.Name, req.Data.Login, req.Data.Password,
			req.Key.Version,
			req.Key.Key,
		)
		err := row.Scan(&result.ID)
		if err != nil {
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("create password error: %w", err)
	}
	return &EncryptedPassword{Item: result, DataKey: req.Key}, nil
}

func (r *SecretRepo) UpdatePassword(ctx context.Context, req model.UpdatePasswordRequest) (*EncryptedPassword, error) {
	result := req.Data
	query := "UPDATE passwords SET name=$1, login=$2, password=$3, sc_version=$4, sc=$5 WHERE id=$6 AND user_id=$7;"

	fun := func() error {
		result, err := r.db.ExecContext(
			ctx, query,
			req.Data.Name, req.Data.Login, req.Data.Password,
			req.Key.Version, req.Key.Key,
			req.Data.ID, req.UserID,
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
		return nil, fmt.Errorf("update password error: %w", err)
	}
	return &EncryptedPassword{Item: result, DataKey: req.Key}, nil
}

func (r *SecretRepo) CreateNote(ctx context.Context, req model.UpdateNoteRequest) (*EncryptedNote, error) {
	result := req.Data
	query := `
	INSERT INTO 
		notes(user_id, name, note, sc_version, sc) 
	values($1, $2, $3, $4, $5)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			req.UserID, req.Data.Name, req.Data.Note,
			req.Key.Version,
			req.Key.Key,
		)
		err := row.Scan(&result.ID)
		if err != nil {
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("create password error: %w", err)
	}
	return &EncryptedNote{Item: result, DataKey: req.Key}, nil
}

func (r *SecretRepo) UpdateNote(ctx context.Context, req model.UpdateNoteRequest) (*EncryptedNote, error) {
	result := req.Data
	query := "UPDATE notes SET name=$1, note=$2, sc_version=$3, sc=$4 WHERE id=$5 AND user_id=$6;"

	fun := func() error {
		result, err := r.db.ExecContext(
			ctx, query,
			req.Data.Name, req.Data.Note,
			req.Key.Version, req.Key.Key,
			req.Data.ID, req.UserID,
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
		return nil, fmt.Errorf("update password error: %w", err)
	}
	return &EncryptedNote{Item: result, DataKey: req.Key}, nil
}

func (r *SecretRepo) CreateCard(ctx context.Context, req model.EncryptedCard, userID int64) (*EncryptedCard, error) {
	result := EncryptedCard{Meta: req.Meta, Payload: req.Payload, DataKey: req.DataKey}
	query := `
	INSERT INTO 
		cards(user_id, name, payload, sc_version, sc) 
	values($1, $2, $3, $4, $5)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			userID, req.Meta.Name, req.Payload,

			req.DataKey.Version,
			req.DataKey.Key,
		)
		err := row.Scan(&result.Meta.ID)
		if err != nil {
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("create password error: %w", err)
	}
	return &result, nil
}

func (r *SecretRepo) UpdateCard(ctx context.Context, req model.EncryptedCard, userID int64) (*EncryptedCard, error) {
	result := EncryptedCard{Meta: req.Meta, Payload: req.Payload, DataKey: req.DataKey}
	query := "UPDATE cards SET name=$1, payload=$2, sc_version=$3, sc=$4 WHERE id=$5 AND user_id=$6;"

	fun := func() error {
		result, err := r.db.ExecContext(
			ctx, query,
			req.Meta.Name, req.Payload,
			req.DataKey.Version, req.DataKey.Key,
			req.Meta.ID, userID,
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
		return nil, fmt.Errorf("update password error: %w", err)
	}
	return &result, nil
}

func (r *SecretRepo) CreateFileMeta(ctx context.Context, req model.CreateFileRequest, path string, key *model.DataKey) (*EncryptedFileMeta, error) {
	result := EncryptedFileMeta{
		Meta: &model.FileMeta{
			SecretMeta: model.SecretMeta{
				Name: req.Name,
				Type: model.SecretTypeFile,
			},
			Path: path,
		},
		DataKey: key,
	}
	query := `
	INSERT INTO 
		files(user_id, name, path, sc_version, sc) 
	values($1, $2, $3, $4, $5)
	RETURNING id;
	`
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx, query,
			req.UserID, req.Name, path,
			key.Version,
			key.Key,
		)
		err := row.Scan(&result.Meta.ID)
		if err != nil {
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("create password error: %w", err)
	}
	return &result, nil
}

func (r *SecretRepo) GetFileMeta(ctx context.Context, req model.SecretRequest) (*EncryptedFileMeta, error) {
	// var result EncryptedFileMeta
	result := EncryptedFileMeta{
		Meta:    &model.FileMeta{SecretMeta: model.SecretMeta{}},
		DataKey: &model.DataKey{},
	}
	query := `
		SELECT 
			id, name, path, sc_version, sc
		FROM 
			public.files
		WHERE id = $1 AND user_id = $2;
	`

	fun := func() error {
		row := r.db.QueryRowContext(ctx, query, req.ID, req.UserID)

		err := row.Scan(
			&result.Meta.ID,
			&result.Meta.Name,
			&result.Meta.Path,
			&result.DataKey.Version,
			&result.DataKey.Key,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("create password error: %w", err)
	}
	return &result, nil
}
