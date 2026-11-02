package repository

import (
	"context"
	"time"
)

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
}
