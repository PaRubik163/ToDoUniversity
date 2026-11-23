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
			COALESCE(s.is_completed, FALSE), s.completed_at, h.position, h.recurrence,
			COALESCE(array_to_json(h.recurrence_days)::text, '[]')
		FROM homework h
		LEFT JOIN homework_status s ON s.homework_id = h.id
		WHERE h.chat_id = $1
		ORDER BY COALESCE(s.is_completed, FALSE),
			CASE WHEN COALESCE(s.is_completed, FALSE) THEN s.completed_at END DESC NULLS LAST,
			CASE WHEN NOT COALESCE(s.is_completed, FALSE) THEN h.position END ASC,
			h.deadline ASC, h.id ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list homework: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Homework, 0)
	for rows.Next() {
		var item domain.Homework
		var daysJSON string
		var recurrenceDaysJSON string
		var completedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Subject, &item.Title, &item.Priority, &item.Note, &item.Deadline, &daysJSON, &item.Completed, &completedAt, &item.Position, &item.Recurrence, &recurrenceDaysJSON); err != nil {
			return nil, fmt.Errorf("scan homework: %w", err)
		}
		if err := json.Unmarshal([]byte(daysJSON), &item.ReminderDays); err != nil {
			return nil, fmt.Errorf("decode reminder days: %w", err)
		}
		if err := json.Unmarshal([]byte(recurrenceDaysJSON), &item.RecurrenceDays); err != nil {
			return nil, fmt.Errorf("decode recurrence days: %w", err)
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
			INSERT INTO homework (chat_id, subject, title, priority, note, deadline, reminder_days, position, recurrence, recurrence_days)
			VALUES ($1, $2, $3, $4, $5, $6, $7::smallint[], COALESCE((SELECT MAX(position)+1 FROM homework WHERE chat_id=$1), 0), $8, $9::smallint[]) RETURNING id
		)
		INSERT INTO homework_status (homework_id, is_completed)
		SELECT id, FALSE FROM new_homework RETURNING homework_id`,
		userID, item.Subject, item.Title, item.Priority, item.Note, item.Deadline, nullableIntArray(item.ReminderDays), item.Recurrence, intArray(item.RecurrenceDays)).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create homework: %w", err)
	}
	return id, nil
}

func (r *HomeworkRepository) Update(ctx context.Context, userID, id int64, item domain.Homework) error {
	result, err := r.db.ExecContext(ctx, `UPDATE homework SET subject = $1, title = $2, priority = $3, note = $4, deadline = $5,
		reminder_days = CASE WHEN $6::smallint[] IS NULL THEN reminder_days ELSE $6::smallint[] END, recurrence=$7, recurrence_days=$8::smallint[]
		WHERE id = $9 AND chat_id = $10`, item.Subject, item.Title, item.Priority, item.Note, item.Deadline, nullableIntArray(item.ReminderDays), item.Recurrence, intArray(item.RecurrenceDays), id, userID)
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin homework status update: %w", err)
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO homework_status (homework_id, is_completed, completed_at)
		SELECT id, $1, CASE WHEN $1 THEN now() ELSE NULL END
		FROM homework WHERE id = $2 AND chat_id = $3
		ON CONFLICT (homework_id) DO UPDATE SET
			is_completed = EXCLUDED.is_completed, completed_at = EXCLUDED.completed_at`, completed, id, userID)
	if err != nil {
		return fmt.Errorf("set homework status: %w", err)
	}
	if err := requireAffected(result); err != nil {
		return err
	}
	if completed {
		_, err = tx.ExecContext(ctx, `INSERT INTO homework (chat_id, subject, title, priority, note, deadline, reminder_days, position, recurrence, recurrence_days, recurrence_source_id)
			SELECT h.chat_id, h.subject, h.title, h.priority, h.note,
				h.deadline + CASE
					WHEN h.recurrence = 'daily' THEN interval '1 day'
					WHEN h.recurrence = 'weekly' THEN interval '1 week'
					WHEN h.recurrence = 'monthly' THEN interval '1 month'
					WHEN h.recurrence = 'weekly' AND cardinality(h.recurrence_days) > 0 THEN make_interval(days => (SELECT MIN(CASE WHEN ((day % 7 - EXTRACT(DOW FROM h.deadline)::int + 7) % 7) = 0 THEN 7 ELSE ((day % 7 - EXTRACT(DOW FROM h.deadline)::int + 7) % 7) END) FROM unnest(h.recurrence_days) AS days(day)))
				END,
				h.reminder_days, (SELECT COALESCE(MAX(position)+1,0) FROM homework WHERE chat_id=h.chat_id), h.recurrence, h.recurrence_days, h.id
			FROM homework h JOIN homework_status s ON s.homework_id=h.id
			WHERE h.id=$1 AND h.chat_id=$2 AND s.is_completed=TRUE AND h.recurrence <> ''
			ON CONFLICT (recurrence_source_id) WHERE recurrence_source_id IS NOT NULL DO NOTHING`, id, userID)
		if err != nil {
			return fmt.Errorf("create next recurring homework: %w", err)
		}
	}
	return tx.Commit()
}

func (r *HomeworkRepository) Move(ctx context.Context, userID, id int64, direction int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin move homework: %w", err)
	}
	defer tx.Rollback()
	var position int64
	err = tx.QueryRowContext(ctx, `SELECT position FROM homework WHERE id=$1 AND chat_id=$2 FOR UPDATE`, id, userID).Scan(&position)
	if err == sql.ErrNoRows {
		return repository.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("find homework to move: %w", err)
	}
	var neighborID int64
	err = tx.QueryRowContext(ctx, `SELECT id FROM homework h
		WHERE h.chat_id=$1 AND h.id<>$2
		AND COALESCE((SELECT is_completed FROM homework_status WHERE homework_id=h.id), FALSE)=FALSE
		AND COALESCE((SELECT is_completed FROM homework_status WHERE homework_id=$2), FALSE)=FALSE
		AND (($3 < 0 AND (h.position, h.id) < ($4, $2)) OR ($3 > 0 AND (h.position, h.id) > ($4, $2)))
		ORDER BY CASE WHEN $3 < 0 THEN h.position END DESC, CASE WHEN $3 > 0 THEN h.position END ASC, h.id
		LIMIT 1 FOR UPDATE`, userID, id, direction, position).Scan(&neighborID)
	if err == sql.ErrNoRows {
		return tx.Commit()
	}
	if err != nil {
		return fmt.Errorf("find neighboring homework: %w", err)
	}
	var neighborPosition int64
	if err := tx.QueryRowContext(ctx, `SELECT position FROM homework WHERE id=$1 FOR UPDATE`, neighborID).Scan(&neighborPosition); err != nil {
		return fmt.Errorf("lock neighboring homework: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE homework SET position=$1 WHERE id=$2`, neighborPosition, id); err != nil {
		return fmt.Errorf("move homework: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE homework SET position=$1 WHERE id=$2`, position, neighborID); err != nil {
		return fmt.Errorf("move neighboring homework: %w", err)
	}
	return tx.Commit()
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

func intArray(values []int) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.Itoa(value)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

var _ repository.HomeworkRepository = (*HomeworkRepository)(nil)
