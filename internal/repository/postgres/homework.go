// Package postgres contains PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"homework-reminder-bot/internal/domain"
	"homework-reminder-bot/internal/repository"
)

type HomeworkRepository struct{ db *sql.DB }

func NewHomeworkRepository(db *sql.DB) *HomeworkRepository {
	return &HomeworkRepository{db: db}
}

func (r *HomeworkRepository) List(ctx context.Context, userID int64) ([]domain.Homework, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT h.id, h.subject, h.title, h.priority, h.note, h.deadline,
			COALESCE(array_to_json(h.reminder_days)::text, '[]'),
			COALESCE(s.is_completed, FALSE), s.completed_at
		FROM homework h
		LEFT JOIN homework_status s ON s.homework_id = h.id
		WHERE h.chat_id = $1
		ORDER BY COALESCE(s.is_completed, FALSE),
			CASE WHEN COALESCE(s.is_completed, FALSE) THEN s.completed_at END DESC NULLS LAST,
			h.deadline ASC, h.id ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list homework: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Homework, 0)
	for rows.Next() {
		var item domain.Homework
		var daysJSON string
		var completedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Subject, &item.Title, &item.Priority, &item.Note, &item.Deadline, &daysJSON, &item.Completed, &completedAt); err != nil {
			return nil, fmt.Errorf("scan homework: %w", err)
		}
		if err := json.Unmarshal([]byte(daysJSON), &item.ReminderDays); err != nil {
			return nil, fmt.Errorf("decode reminder days: %w", err)
		}
		if completedAt.Valid {
			value := completedAt.Time
			item.CompletedAt = &value
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate homework: %w", err)
	}
	return items, nil
}

func (r *HomeworkRepository) Create(ctx context.Context, userID int64, item domain.Homework) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		WITH new_homework AS (
			INSERT INTO homework (chat_id, subject, title, priority, note, deadline, reminder_days)
			VALUES ($1, $2, $3, $4, $5, $6, $7::smallint[]) RETURNING id
		)
		INSERT INTO homework_status (homework_id, is_completed)
		SELECT id, FALSE FROM new_homework RETURNING homework_id`,
		userID, item.Subject, item.Title, item.Priority, item.Note, item.Deadline, nullableIntArray(item.ReminderDays)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create homework: %w", err)
	}
	return id, nil
}

func (r *HomeworkRepository) Update(ctx context.Context, userID, id int64, item domain.Homework) error {
	result, err := r.db.ExecContext(ctx, `UPDATE homework SET subject = $1, title = $2, priority = $3, note = $4, deadline = $5,
		reminder_days = CASE WHEN $6::smallint[] IS NULL THEN reminder_days ELSE $6::smallint[] END
		WHERE id = $7 AND chat_id = $8`, item.Subject, item.Title, item.Priority, item.Note, item.Deadline, nullableIntArray(item.ReminderDays), id, userID)
	if err != nil {
		return fmt.Errorf("update homework: %w", err)
	}
	return requireAffected(result)
}

func (r *HomeworkRepository) Delete(ctx context.Context, userID, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM homework WHERE id = $1 AND chat_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete homework: %w", err)
	}
	return requireAffected(result)
}

func (r *HomeworkRepository) SetCompleted(ctx context.Context, userID, id int64, completed bool) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO homework_status (homework_id, is_completed, completed_at)
		SELECT id, $1, CASE WHEN $1 THEN now() ELSE NULL END
		FROM homework WHERE id = $2 AND chat_id = $3
		ON CONFLICT (homework_id) DO UPDATE SET
			is_completed = EXCLUDED.is_completed, completed_at = EXCLUDED.completed_at`, completed, id, userID)
	if err != nil {
		return fmt.Errorf("set homework status: %w", err)
	}
	return requireAffected(result)
}

func requireAffected(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected row count: %w", err)
	}
	if count == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func nullableIntArray(values []int) any {
	if len(values) == 0 {
		return nil
	}
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.Itoa(value)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

var _ repository.HomeworkRepository = (*HomeworkRepository)(nil)
