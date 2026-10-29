// Package handler adapts HTTP requests and responses to application services.
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"homework-reminder-bot/internal/domain"
	"homework-reminder-bot/internal/repository"
	"homework-reminder-bot/internal/service"
)

type HomeworkHandler struct {
	service  service.HomeworkService
	timezone *time.Location
}

func NewHomeworkHandler(homeworkService service.HomeworkService, timezone *time.Location) *HomeworkHandler {
	return &HomeworkHandler{service: homeworkService, timezone: timezone}
}

type homeworkRequest struct {
	Subject      string `json:"subject"`
	Title        string `json:"title"`
	Priority     string `json:"priority"`
	Note         string `json:"note"`
	Deadline     string `json:"deadline"`
	ReminderDays []int  `json:"reminder_days"`
}

type homeworkResponse struct {
	ID           int64  `json:"id"`
	Subject      string `json:"subject"`
	Title        string `json:"title"`
	Priority     string `json:"priority"`
	Note         string `json:"note"`
	Deadline     string `json:"deadline"`
	ReminderDays []int  `json:"reminder_days"`
	DaysLeft     int    `json:"days_left"`
	Completed    bool   `json:"completed"`
	CompletedAt  string `json:"completed_at,omitempty"`
}

func (h *HomeworkHandler) Collection(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.service.List(r.Context(), userID)
		if err != nil {
			h.internalError(w, "list homework", err)
			return
		}
		response := make([]homeworkResponse, 0, len(items))
		for _, item := range items {
			response = append(response, h.present(item))
		}
		writeJSON(w, response)
	case http.MethodPost:
		var req homeworkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		item, err := h.service.Create(r.Context(), userID, toInput(req))
		if err != nil {
			h.writeServiceError(w, "create homework", err)
			return
		}
		writeJSON(w, h.present(item))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HomeworkHandler) Item(w http.ResponseWriter, r *http.Request, userID int64) {
	path := strings.TrimPrefix(r.URL.Path, "/api/homework/")
	completedPath := strings.HasSuffix(path, "/status")
	if completedPath {
		path = strings.TrimSuffix(path, "/status")
	}
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if completedPath {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Completed bool `json:"completed"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.service.SetCompleted(r.Context(), userID, id, req.Completed); err != nil {
			h.writeServiceError(w, "set homework status", err)
			return
		}
		writeJSON(w, map[string]bool{"completed": req.Completed})
		return
	}
	switch r.Method {
	case http.MethodPut:
		var req homeworkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		item, err := h.service.Update(r.Context(), userID, id, toInput(req))
		if err != nil {
			h.writeServiceError(w, "update homework", err)
			return
		}
		writeJSON(w, h.present(item))
	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), userID, id); err != nil {
			h.writeServiceError(w, "delete homework", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HomeworkHandler) present(item domain.Homework) homeworkResponse {
	now := time.Now().In(h.timezone)
	deadline := item.Deadline.In(h.timezone)
	response := homeworkResponse{
		ID: item.ID, Subject: item.Subject, Title: item.Title, Priority: item.Priority, Note: item.Note,
		Deadline: deadline.Format("2006-01-02"), ReminderDays: item.ReminderDays,
		DaysLeft: domain.CalendarDays(now, deadline, h.timezone), Completed: item.Completed,
	}
	if item.CompletedAt != nil {
		response.CompletedAt = item.CompletedAt.In(h.timezone).Format("2006-01-02")
	}
	return response
}

func (h *HomeworkHandler) writeServiceError(w http.ResponseWriter, operation string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, strings.TrimPrefix(err.Error(), service.ErrInvalidInput.Error()+": "), http.StatusBadRequest)
	case errors.Is(err, repository.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	default:
		h.internalError(w, operation, err)
	}
}

func (h *HomeworkHandler) internalError(w http.ResponseWriter, operation string, err error) {
	log.Printf("%s: %v", operation, err)
	http.Error(w, "database error", http.StatusInternalServerError)
}

func toInput(req homeworkRequest) service.HomeworkInput {
	return service.HomeworkInput{Subject: req.Subject, Title: req.Title, Priority: req.Priority, Note: req.Note, Deadline: req.Deadline, ReminderDays: req.ReminderDays}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
