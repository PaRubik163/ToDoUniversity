package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed miniapp/*
var miniAppFiles embed.FS

// apiHomework is the JSON shape sent to the Mini App frontend.
type apiHomework struct {
	ID       int64  `json:"id"`
	Subject  string `json:"subject"`
	Title    string `json:"title"`
	Deadline string `json:"deadline"` // YYYY-MM-DD
	DaysLeft int    `json:"days_left"`
}

type createHomeworkRequest struct {
	Subject  string `json:"subject"`
	Title    string `json:"title"`
	Deadline string `json:"deadline"` // YYYY-MM-DD, matches <input type="date">
}

type telegramUser struct {
	ID int64 `json:"id"`
}

// startWebServer serves the Mini App static files and its JSON API. It blocks
// until the context is cancelled or the server fails to start.
func (b *Bot) startWebServer(ctx context.Context, addr string) error {
	sub, err := fs.Sub(miniAppFiles, "miniapp")
	if err != nil {
		return fmt.Errorf("mount miniapp assets: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/homework", b.withAuth(b.handleHomeworkCollection))
	mux.HandleFunc("/api/homework/", b.withAuth(b.handleHomeworkItem))

	server := &http.Server{
		Addr:              addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("mini app server listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s", r.Method, r.URL.Path)
	})
}

// withAuth validates the Telegram Mini App initData sent by the frontend in the
// X-Telegram-Init-Data header and passes the authenticated Telegram user id
// through to the wrapped handler. The user id doubles as the chat_id used
// everywhere else in this app, since the bot only ever talks to users 1:1.
func (b *Bot) withAuth(next func(w http.ResponseWriter, r *http.Request, userID int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initData := r.Header.Get("X-Telegram-Init-Data")
		if initData == "" {
			http.Error(w, "missing init data", http.StatusUnauthorized)
			return
		}
		userID, err := validateInitData(initData, b.config.Token)
		if err != nil {
			http.Error(w, "invalid init data: "+err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r, userID)
	}
}

// validateInitData checks the Telegram WebApp init-data signature per
// https://core.telegram.org/bots/webapps#validating-data-received-via-the-mini-app
// and returns the authenticated user's Telegram id.
func validateInitData(initData, botToken string) (int64, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, fmt.Errorf("parse init data: %w", err)
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return 0, errors.New("missing hash")
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+values.Get(k))
	}
	dataCheckString := strings.Join(pairs, "\n")

	secretKeyMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyMAC.Write([]byte(botToken))
	secretKey := secretKeyMAC.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(computedHash), []byte(receivedHash)) {
		return 0, errors.New("hash mismatch")
	}

	if authDateStr := values.Get("auth_date"); authDateStr != "" {
		if authDate, err := strconv.ParseInt(authDateStr, 10, 64); err == nil {
			if time.Since(time.Unix(authDate, 0)) > 24*time.Hour {
				return 0, errors.New("init data expired")
			}
		}
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return 0, errors.New("missing user field")
	}
	var user telegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return 0, fmt.Errorf("parse user field: %w", err)
	}
	if user.ID == 0 {
		return 0, errors.New("missing user id")
	}
	return user.ID, nil
}

func (b *Bot) handleHomeworkCollection(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodGet:
		b.apiListHomework(w, r, userID)
	case http.MethodPost:
		b.apiCreateHomework(w, r, userID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (b *Bot) handleHomeworkItem(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/homework/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	result, err := b.db.ExecContext(r.Context(), "DELETE FROM homework WHERE id = $1 AND chat_id = $2", id, userID)
	if err != nil {
		log.Printf("delete homework: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("delete homework rows affected: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (b *Bot) apiListHomework(w http.ResponseWriter, r *http.Request, userID int64) {
	rows, err := b.db.QueryContext(r.Context(), `SELECT id, subject, title, deadline FROM homework WHERE chat_id = $1 ORDER BY deadline ASC, id ASC`, userID)
	if err != nil {
		log.Printf("list homework: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	now := time.Now().In(b.config.TZ)
	items := make([]apiHomework, 0)
	for rows.Next() {
		var item homework
		if err := rows.Scan(&item.ID, &item.Subject, &item.Title, &item.Deadline); err != nil {
			log.Printf("scan homework: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		deadline := item.Deadline.In(b.config.TZ)
		items = append(items, apiHomework{
			ID:       item.ID,
			Subject:  item.Subject,
			Title:    item.Title,
			Deadline: deadline.Format("2006-01-02"),
			DaysLeft: calendarDays(now, deadline, b.config.TZ),
		})
	}
	if err := rows.Err(); err != nil {
		log.Printf("list homework rows: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}

func (b *Bot) apiCreateHomework(w http.ResponseWriter, r *http.Request, userID int64) {
	var req createHomeworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Subject = strings.TrimSpace(req.Subject)
	req.Title = strings.TrimSpace(req.Title)
	if req.Subject == "" || req.Title == "" || req.Deadline == "" {
		http.Error(w, "subject, title and deadline are required", http.StatusBadRequest)
		return
	}
	deadline, err := time.ParseInLocation("2006-01-02", req.Deadline, b.config.TZ)
	if err != nil {
		http.Error(w, "deadline must be in YYYY-MM-DD format", http.StatusBadRequest)
		return
	}

	var id int64
	err = b.db.QueryRowContext(r.Context(),
		`INSERT INTO homework (chat_id, subject, title, deadline) VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, req.Subject, req.Title, deadline,
	).Scan(&id)
	if err != nil {
		log.Printf("create homework: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	now := time.Now().In(b.config.TZ)
	writeJSON(w, apiHomework{
		ID:       id,
		Subject:  req.Subject,
		Title:    req.Title,
		Deadline: deadline.Format("2006-01-02"),
		DaysLeft: calendarDays(now, deadline, b.config.TZ),
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}

// setMenuButton configures the bot's persistent menu button to open the Mini
// App. It's a no-op if MiniAppURL isn't configured.
func (b *Bot) setMenuButton(ctx context.Context) error {
	if b.config.MiniAppURL == "" {
		return nil
	}
	type webAppInfo struct {
		URL string `json:"url"`
	}
	type menuButton struct {
		Type   string      `json:"type"`
		Text   string      `json:"text"`
		WebApp *webAppInfo `json:"web_app"`
	}
	payload := struct {
		MenuButton menuButton `json:"menu_button"`
	}{
		MenuButton: menuButton{Type: "web_app", Text: "Задания", WebApp: &webAppInfo{URL: b.config.MiniAppURL}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/setChatMenuButton", b.config.Token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("telegram status: %s", resp.Status)
	}
	return nil
}

// sendWebAppButton sends a text message with an inline button that opens the
// Mini App directly inside Telegram.
func (b *Bot) sendWebAppButton(ctx context.Context, chatID int64, text, buttonText, appURL string) error {
	type webAppInfo struct {
		URL string `json:"url"`
	}
	type inlineButton struct {
		Text   string      `json:"text"`
		WebApp *webAppInfo `json:"web_app"`
	}
	keyboard := struct {
		InlineKeyboard [][]inlineButton `json:"inline_keyboard"`
	}{
		InlineKeyboard: [][]inlineButton{{{Text: buttonText, WebApp: &webAppInfo{URL: appURL}}}},
	}
	markup, err := json.Marshal(keyboard)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.config.Token)
	form := url.Values{"chat_id": {fmt.Sprint(chatID)}, "text": {text}, "reply_markup": {string(markup)}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("telegram status: %s", resp.Status)
	}
	return nil
}
