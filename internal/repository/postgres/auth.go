package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"homework-reminder-bot/internal/repository"
)

type AuthRepository struct{ db *sql.DB }

func NewAuthRepository(db *sql.DB) *AuthRepository { return &AuthRepository{db: db} }

func (r *AuthRepository) Create(ctx context.Context, email, passwordHash string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var id int64
	if err = tx.QueryRowContext(ctx, `SELECT nextval('website_user_ids')`).Scan(&id); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO auth_accounts(id,email,password_hash) VALUES($1,$2,$3)`, id, email, passwordHash); err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO user_settings(chat_id,language) VALUES($1,'ru')`, id); err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}
func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (repository.AuthAccount, error) {
	var a repository.AuthAccount
	err := r.db.QueryRowContext(ctx, `SELECT id,email,password_hash FROM auth_accounts WHERE email=$1`, email).Scan(&a.ID, &a.Email, &a.PasswordHash)
	return a, err
}
func (r *AuthRepository) SaveRefresh(ctx context.Context, hash []byte, userID int64, expires time.Time) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO auth_refresh_tokens(token_hash,user_id,expires_at) VALUES($1,$2,$3)`, hash, userID, expires)
	return err
}
func (r *AuthRepository) RotateRefresh(ctx context.Context, oldHash, newHash []byte, expires time.Time) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var id int64
	err = tx.QueryRowContext(ctx, `DELETE FROM auth_refresh_tokens WHERE token_hash=$1 AND expires_at>now() RETURNING user_id`, oldHash).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO auth_refresh_tokens(token_hash,user_id,expires_at) VALUES($1,$2,$3)`, newHash, id, expires); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}
func (r *AuthRepository) RevokeRefresh(ctx context.Context, hash []byte) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth_refresh_tokens WHERE token_hash=$1`, hash)
	return err
}

var _ repository.AuthRepository = (*AuthRepository)(nil)
