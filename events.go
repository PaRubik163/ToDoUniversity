package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type calendarEvent struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	EventAt       string `json:"event_at"`
	ReminderHours int    `json:"reminder_hours"`
	Completed     bool   `json:"completed"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

type calendarEventRequest struct {
	Title         string `json:"title"`
	EventAt       string `json:"event_at"`
	ReminderHours int    `json:"reminder_hours"`
}

func (b *Bot) handleEvents(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodGet:
		month := strings.TrimSpace(r.URL.Query().Get("month"))
		start, err := time.ParseInLocation("2006-01", month, b.config.Timezone)
		if err != nil || start.Format("2006-01") != month {
			http.Error(w, "month must be in YYYY-MM format", http.StatusBadRequest)
			return
		}
		rows, err := b.db.QueryContext(r.Context(), `SELECT id, title, event_at, reminder_hours, is_completed, completed_at FROM calendar_events WHERE chat_id = $1 AND event_at >= $2 AND event_at < $3 ORDER BY event_at`, userID, start, start.AddDate(0, 1, 0))
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		items := make([]calendarEvent, 0)
		for rows.Next() {
			var item calendarEvent
			var at time.Time
			var completedAt *time.Time
			if err := rows.Scan(&item.ID, &item.Title, &at, &item.ReminderHours, &item.Completed, &completedAt); err != nil {
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
			item.EventAt = at.In(b.config.Timezone).Format(time.RFC3339)
			if completedAt != nil {
				item.CompletedAt = completedAt.In(b.config.Timezone).Format(time.RFC3339)
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, items)
	case http.MethodPost:
		var req calendarEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		req.Title = strings.TrimSpace(req.Title)
		at, err := time.ParseInLocation("2006-01-02T15:04", strings.TrimSpace(req.EventAt), b.config.Timezone)
		if req.Title == "" || len(req.Title) > 300 || err != nil || req.ReminderHours < 0 || req.ReminderHours > 168 {
			http.Error(w, "title, event_at, or reminder_hours is invalid", http.StatusBadRequest)
			return
		}
		var id int64
		if err := b.db.QueryRowContext(r.Context(), `INSERT INTO calendar_events (chat_id, title, event_at, reminder_hours) VALUES ($1,$2,$3,$4) RETURNING id`, userID, req.Title, at, req.ReminderHours).Scan(&id); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, calendarEvent{ID: id, Title: req.Title, EventAt: at.Format(time.RFC3339), ReminderHours: req.ReminderHours})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (b *Bot) handleEventItem(w http.ResponseWriter, r *http.Request, userID int64) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/"), "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}
	switch {
	case r.Method == http.MethodPut && len(parts) == 2 && parts[1] == "status":
		var req struct {
			Completed bool `json:"completed"`
		}
		if json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := b.db.ExecContext(r.Context(), `UPDATE calendar_events SET is_completed=$3, completed_at=CASE WHEN $3 THEN now() ELSE NULL END WHERE id=$1 AND chat_id=$2`, id, userID, req.Completed)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if n, _ := result.RowsAffected(); n == 0 {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete && len(parts) == 1:
		result, err := b.db.ExecContext(r.Context(), `DELETE FROM calendar_events WHERE id=$1 AND chat_id=$2`, id, userID)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if n, _ := result.RowsAffected(); n == 0 {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
