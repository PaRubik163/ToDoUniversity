// Package repository defines persistence contracts used by application
// services.
package repository

import (
	"context"
	"errors"

	"homework-reminder-bot/internal/domain"
)

var ErrNotFound = errors.New("homework not found")

// HomeworkRepository is the storage contract for task operations.
type HomeworkRepository interface {
	List(ctx context.Context, userID int64) ([]domain.Homework, error)
	Create(ctx context.Context, userID int64, item domain.Homework) (int64, error)
	Update(ctx context.Context, userID, id int64, item domain.Homework) error
	Delete(ctx context.Context, userID, id int64) error
	SetCompleted(ctx context.Context, userID, id int64, completed bool) error
}
