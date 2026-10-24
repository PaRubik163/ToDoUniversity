package repository

import "context"

// UserRepository owns persistence operations for app users.
type UserRepository interface {
	EnsureExists(ctx context.Context, userID int64, language string) error
}
