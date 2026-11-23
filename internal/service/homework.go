// Package service contains application use cases.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"homework-reminder-bot/internal/domain"
	"homework-reminder-bot/internal/repository"
)

var ErrInvalidInput = errors.New("invalid homework input")

type HomeworkInput struct {
	Subject        string
	Title          string
	Priority       string
	Note           string
	Deadline       string
	ReminderDays   []int
	Recurrence     string
	RecurrenceDays []int
}

type HomeworkService interface {
	List(ctx context.Context, userID int64) ([]domain.Homework, error)
	Create(ctx context.Context, userID int64, input HomeworkInput) (domain.Homework, error)
	Update(ctx context.Context, userID, id int64, input HomeworkInput) (domain.Homework, error)
	Delete(ctx context.Context, userID, id int64) error
	SetCompleted(ctx context.Context, userID, id int64, completed bool) error
	Move(ctx context.Context, userID, id int64, direction int) error
}

type homeworkService struct {
	repository repository.HomeworkRepository
	timezone   *time.Location
}

func NewHomeworkService(repo repository.HomeworkRepository, timezone *time.Location) HomeworkService {
	return &homeworkService{repository: repo, timezone: timezone}
}

func (s *homeworkService) List(ctx context.Context, userID int64) ([]domain.Homework, error) {
	return s.repository.List(ctx, userID)
}

func (s *homeworkService) Create(ctx context.Context, userID int64, input HomeworkInput) (domain.Homework, error) {
	item, err := s.validated(input)
	if err != nil {
		return domain.Homework{}, err
	}
	item.ID, err = s.repository.Create(ctx, userID, item)
	if err != nil {
		return domain.Homework{}, err
	}
	return item, nil
}

func (s *homeworkService) Update(ctx context.Context, userID, id int64, input HomeworkInput) (domain.Homework, error) {
	item, err := s.validated(input)
	if err != nil {
		return domain.Homework{}, err
	}
	if err := s.repository.Update(ctx, userID, id, item); err != nil {
		return domain.Homework{}, err
	}
	item.ID = id
	return item, nil
}

func (s *homeworkService) Delete(ctx context.Context, userID, id int64) error {
	return s.repository.Delete(ctx, userID, id)
}

func (s *homeworkService) SetCompleted(ctx context.Context, userID, id int64, completed bool) error {
	return s.repository.SetCompleted(ctx, userID, id, completed)
}

func (s *homeworkService) Move(ctx context.Context, userID, id int64, direction int) error {
	if direction != -1 && direction != 1 {
		return fmt.Errorf("%w: direction must be -1 or 1", ErrInvalidInput)
	}
	return s.repository.Move(ctx, userID, id, direction)
}

func (s *homeworkService) validated(input HomeworkInput) (domain.Homework, error) {
	input.Subject, input.Title = strings.TrimSpace(input.Subject), strings.TrimSpace(input.Title)
	if input.Title == "" || input.Deadline == "" {
		return domain.Homework{}, fmt.Errorf("%w: title and deadline are required", ErrInvalidInput)
	}
	if len(input.Subject) > 200 || len(input.Title) > 300 || len(input.Note) > 2000 {
		return domain.Homework{}, fmt.Errorf("%w: subject/title/note is too long", ErrInvalidInput)
	}
	input.Priority = strings.TrimSpace(input.Priority)
	if input.Priority == "" {
		input.Priority = "medium"
	}
	if input.Priority != "low" && input.Priority != "medium" && input.Priority != "high" {
		return domain.Homework{}, fmt.Errorf("%w: priority must be low, medium or high", ErrInvalidInput)
	}
	input.Recurrence = strings.TrimSpace(input.Recurrence)
	if input.Recurrence != "" && input.Recurrence != "daily" && input.Recurrence != "weekly" && input.Recurrence != "monthly" {
		return domain.Homework{}, fmt.Errorf("%w: recurrence value is invalid", ErrInvalidInput)
	}
	if input.Recurrence != "weekly" && len(input.RecurrenceDays) > 0 {
		return domain.Homework{}, fmt.Errorf("%w: recurrence days require weekly recurrence", ErrInvalidInput)
	}
	seenDays := make(map[int]bool, len(input.RecurrenceDays))
	for _, day := range input.RecurrenceDays {
		if day < 1 || day > 7 || seenDays[day] {
			return domain.Homework{}, fmt.Errorf("%w: recurrence days must be unique values from 1 to 7", ErrInvalidInput)
		}
		seenDays[day] = true
	}
	if input.Recurrence == "weekly" && len(input.RecurrenceDays) == 0 {
		return domain.Homework{}, fmt.Errorf("%w: choose at least one weekday", ErrInvalidInput)
	}
	input.Note = strings.TrimSpace(input.Note)
	if len(input.ReminderDays) > 0 && (len(input.ReminderDays) != 1 || input.ReminderDays[0] < 1 || input.ReminderDays[0] > 20) {
		return domain.Homework{}, fmt.Errorf("%w: reminder_days must contain one value from 1 to 20", ErrInvalidInput)
	}
	deadline, err := time.ParseInLocation("2006-01-02", input.Deadline, s.timezone)
	if err != nil || deadline.Format("2006-01-02") != input.Deadline {
		return domain.Homework{}, fmt.Errorf("%w: deadline must use YYYY-MM-DD format", ErrInvalidInput)
	}
	return domain.Homework{Subject: input.Subject, Title: input.Title, Priority: input.Priority, Note: input.Note, Deadline: deadline, ReminderDays: input.ReminderDays, Recurrence: input.Recurrence, RecurrenceDays: input.RecurrenceDays}, nil
}
