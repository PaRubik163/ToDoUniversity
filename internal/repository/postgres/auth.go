package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return 0, repository.ErrEmailExists
		}
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
	if _, err := r.db.ExecContext(ctx, `DELETE FROM auth_refresh_tokens WHERE expires_at<=now()`); err != nil {
		return err
	}
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

func (r *AuthRepository) CreateLinkCode(ctx context.Context, hash []byte, accountID int64, expires time.Time) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM auth_link_codes WHERE expires_at<=now()`); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO auth_link_codes(token_hash,account_id,expires_at) VALUES($1,$2,$3)`, hash, accountID, expires)
	return err
}

func (r *AuthRepository) LinkTelegramAccount(ctx context.Context, hash []byte, telegramID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var accountID int64
	if err = tx.QueryRowContext(ctx, `SELECT account_id FROM auth_link_codes WHERE token_hash=$1 AND expires_at>now() FOR UPDATE`, hash).Scan(&accountID); err != nil {
		return err
	}
	var alreadyLinked bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth_accounts WHERE id=$1)`, telegramID).Scan(&alreadyLinked); err != nil {
		return err
	}
	if alreadyLinked {
		return errors.New("Telegram account is already linked")
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM user_settings WHERE chat_id=$1`, accountID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE homework SET chat_id=$1 WHERE chat_id=$2`, telegramID, accountID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE calendar_events SET chat_id=$1 WHERE chat_id=$2`, telegramID, accountID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM auth_link_codes WHERE account_id=$1`, accountID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE auth_accounts SET id=$1 WHERE id=$2`, telegramID, accountID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM auth_link_codes WHERE token_hash=$1`, hash); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *AuthRepository) AccountExists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM auth_accounts WHERE id=$1)`, id).Scan(&exists)
	return exists, err
}

var _ repository.AuthRepository = (*AuthRepository)(nil)
