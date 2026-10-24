package service

import (
	"context"

	"homework-reminder-bot/internal/repository"
)

type UserService struct{ users repository.UserRepository }

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) EnsureExists(ctx context.Context, userID int64, language string) error {
	return s.users.EnsureExists(ctx, userID, language)
}
