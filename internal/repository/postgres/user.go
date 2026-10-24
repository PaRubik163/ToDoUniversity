package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"homework-reminder-bot/internal/repository"
)

type UserRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) EnsureExists(ctx context.Context, userID int64, language string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_settings (chat_id, language)
		VALUES ($1, $2)
		ON CONFLICT (chat_id) DO NOTHING`, userID, language)
	if err != nil {
		return fmt.Errorf("ensure user settings: %w", err)
	}
	return nil
}

var _ repository.UserRepository = (*UserRepository)(nil)
