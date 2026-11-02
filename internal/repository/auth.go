package repository

import (
	"context"
	"errors"
	"time"
)

var ErrEmailExists = errors.New("email already exists")

type AuthAccount struct {
	ID                  int64
	Email, PasswordHash string
}

type AuthRepository interface {
	Create(context.Context, string, string) (int64, error)
	FindByEmail(context.Context, string) (AuthAccount, error)
	SaveRefresh(context.Context, []byte, int64, time.Time) error
	RotateRefresh(context.Context, []byte, []byte, time.Time) (int64, error)
	RevokeRefresh(context.Context, []byte) error
	CreateLinkCode(context.Context, []byte, int64, time.Time) error
	LinkTelegramAccount(context.Context, []byte, int64) error
	AccountExists(context.Context, int64) (bool, error)
}
