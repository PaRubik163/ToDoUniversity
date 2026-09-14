package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const pollTimeout = 25 * time.Second

type Config struct {
	Token      string
	TZ         *time.Location
	MiniAppURL string
}

type Bot struct {
	config Config
	db     *sql.DB
	client *http.Client
}

type updateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int64 `json:"message_id"`
	Chat      struct {
		ID int64 `json:"id"`
	} `json:"chat"`
	Text string `json:"text"`
}

type homework struct {
	ID       int64
	Subject  string
	Title    string
	Deadline time.Time
}

func main() {
	ctx := context.Background()
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if token == "" || databaseURL == "" {
		log.Fatal("set TELEGRAM_BOT_TOKEN and DATABASE_URL")
	}

	tzName := os.Getenv("BOT_TIMEZONE")
	if tzName == "" {
		tzName = "Europe/Moscow"
	}
	tz, err := time.LoadLocation(tzName)
	if err != nil {
		log.Fatalf("invalid BOT_TIMEZONE: %v", err)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("postgres is unavailable: %v", err)
	}
	if err := migrate(ctx, db); err != nil {
		log.Fatal(err)
	}

	miniAppURL := strings.TrimSpace(os.Getenv("MINI_APP_URL"))
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	b := &Bot{
		config: Config{Token: token, TZ: tz, MiniAppURL: miniAppURL},
		db:     db,
		client: &http.Client{Timeout: pollTimeout + 10*time.Second},
	}

	if miniAppURL != "" {
		if err := b.setMenuButton(ctx); err != nil {
			log.Printf("set menu button: %v", err)
		}
	} else {
		log.Println("MINI_APP_URL not set: skipping Telegram menu button, mini app will only be reachable directly by URL")
	}

	go b.reminderLoop(ctx)
	go func() {
		if err := b.startWebServer(ctx, ":"+port); err != nil {
			log.Fatalf("web server: %v", err)
		}
	}()

	log.Println("bot is running")
	if err := b.poll(ctx); err != nil {
		log.Fatal(err)
	}
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

// migrate applies each new SQL file exactly once. Add future schema changes as
// migrations/002_description.sql, migrations/003_description.sql, and so on.
func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		return fmt.Errorf("create migration registry: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	seen := make(map[int]bool)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		prefix, _, ok := strings.Cut(entry.Name(), "_")
		if !ok {
			return fmt.Errorf("migration %q must start with a numeric version and underscore", entry.Name())
		}
		version, err := strconv.Atoi(prefix)
		if err != nil || version < 1 || seen[version] {
			return fmt.Errorf("invalid or duplicate migration version in %q", entry.Name())
		}
		seen[version] = true

		sqlBytes, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read migration %q: %w", entry.Name(), err)
		}
		if err := applyMigration(ctx, db, version, string(sqlBytes)); err != nil {
			return fmt.Errorf("apply migration %q: %w", entry.Name(), err)
		}
	}
	return nil
}

func applyMigration(ctx context.Context, db *sql.DB, version int, sqlText string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var applied bool
	if err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&applied); err != nil {
		return err
	}
	if applied {
		return tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *Bot) poll(ctx context.Context) error {
	var offset int64
	for {
		updates, err := b.getUpdates(ctx, offset)
		if err != nil {
			log.Printf("getUpdates: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message == nil || strings.TrimSpace(update.Message.Text) == "" {
				continue
			}
			if err := b.handleMessage(ctx, update.Message); err != nil {
				log.Printf("handle message: %v", err)
				_ = b.sendMessage(ctx, update.Message.Chat.ID, "Не получилось обработать сообщение. Попробуйте ещё раз.")
			}
		}
	}
}

func (b *Bot) getUpdates(ctx context.Context, offset int64) ([]Update, error) {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", b.config.Token)
	form := url.Values{"timeout": {"25"}, "offset": {fmt.Sprint(offset)}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result updateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, errors.New("telegram returned an error")
	}
	return result.Result, nil
}

func (b *Bot) handleMessage(ctx context.Context, msg *Message) error {
	text := strings.TrimSpace(msg.Text)
	switch {
	case text == "/start" || text == "/help":
		if b.config.MiniAppURL != "" {
			err := b.sendWebAppButton(ctx, msg.Chat.ID, helpText, "📋 Открыть список заданий", b.config.MiniAppURL)
			if err == nil {
				return nil
			}
			log.Printf("send web app button: %v", err)
		}
		return b.sendMessage(ctx, msg.Chat.ID, helpText)
	case text == "/list" || text == "Список":
		return b.sendHomeworkList(ctx, msg.Chat.ID)
	case strings.HasPrefix(text, "/delete") || strings.HasPrefix(text, "Удалить"):
		return b.deleteHomework(ctx, msg.Chat.ID, text)
	case strings.HasPrefix(text, "/edit"):
		return b.editHomeworkText(ctx, msg.Chat.ID, text)
	case strings.HasPrefix(text, "/remindtime"):
		return b.setReminderTimeText(ctx, msg.Chat.ID, text)
	case strings.Contains(strings.ToLower(text), "предмет:"):
		subject, title, deadline, err := parseHomework(text, b.config.TZ, time.Now().In(b.config.TZ))
		if err != nil {
			return b.sendMessage(ctx, msg.Chat.ID, "Не удалось прочитать задание: "+err.Error()+"\n\n"+formatExample)
		}
		var id int64
		err = b.db.QueryRowContext(ctx, `INSERT INTO homework (chat_id, subject, title, deadline) VALUES ($1, $2, $3, $4) RETURNING id`, msg.Chat.ID, subject, title, deadline).Scan(&id)
		if err != nil {
			return err
		}
		return b.sendMessage(ctx, msg.Chat.ID, fmt.Sprintf("Записала задание №%d. Напомню за 7 дней до дедлайна и далее каждый день.", id))
	default:
		return b.sendMessage(ctx, msg.Chat.ID, "Я понимаю добавление задания, /list, /edit, /delete и /remindtime.\n\n"+formatExample)
	}
}

func (b *Bot) deleteHomework(ctx context.Context, chatID int64, text string) error {
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return b.sendMessage(ctx, chatID, "Укажите номер задания, например: /delete 3")
	}
	var id int64
	if _, err := fmt.Sscan(parts[1], &id); err != nil || id < 1 {
		return b.sendMessage(ctx, chatID, "Номер задания должен быть числом: /delete 3")
	}
	result, err := b.db.ExecContext(ctx, "DELETE FROM homework WHERE id = $1 AND chat_id = $2", id, chatID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return b.sendMessage(ctx, chatID, "Задание с таким номером не найдено.")
	}
	return b.sendMessage(ctx, chatID, "Задание удалено ✅")
}

func (b *Bot) editHomeworkText(ctx context.Context, chatID int64, text string) error {
	lines := strings.SplitN(text, "\n", 2)
	header := strings.Fields(lines[0])
	if len(header) != 2 {
		return b.sendMessage(ctx, chatID, "Формат: /edit НОМЕР, затем поля как при создании.\n\n"+editExample)
	}
	id, err := strconv.ParseInt(header[1], 10, 64)
	if err != nil || id < 1 {
		return b.sendMessage(ctx, chatID, "Номер задания должен быть числом: /edit 3")
	}
	if len(lines) < 2 {
		return b.sendMessage(ctx, chatID, "Не хватает полей после номера.\n\n"+editExample)
	}
	subject, title, deadline, err := parseHomework(lines[1], b.config.TZ, time.Now().In(b.config.TZ))
	if err != nil {
		return b.sendMessage(ctx, chatID, "Не удалось прочитать задание: "+err.Error()+"\n\n"+editExample)
	}
	result, err := b.db.ExecContext(ctx, `UPDATE homework SET subject = $1, title = $2, deadline = $3 WHERE id = $4 AND chat_id = $5`, subject, title, deadline, id, chatID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return b.sendMessage(ctx, chatID, "Задание с таким номером не найдено.")
	}
	return b.sendMessage(ctx, chatID, fmt.Sprintf("Задание №%d обновлено ✅", id))
}

func (b *Bot) setReminderTimeText(ctx context.Context, chatID int64, text string) error {
	parts := strings.Fields(text)
	if len(parts) != 2 {
		return b.sendMessage(ctx, chatID, "Формат: /remindtime ЧЧ:ММ, например /remindtime 19:00")
	}
	hour, minute, err := parseClockTime(parts[1])
	if err != nil {
		return b.sendMessage(ctx, chatID, "Не удалось прочитать время: "+err.Error()+"\n\nФормат: /remindtime 19:00")
	}
	if err := b.setReminderTime(ctx, chatID, hour, minute); err != nil {
		return err
	}
	return b.sendMessage(ctx, chatID, fmt.Sprintf("Готово, буду напоминать каждый день в %02d:%02d.", hour, minute))
}

func (b *Bot) setReminderTime(ctx context.Context, chatID int64, hour, minute int) error {
	_, err := b.db.ExecContext(ctx, `
		INSERT INTO user_settings (chat_id, reminder_hour, reminder_minute, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (chat_id) DO UPDATE SET reminder_hour = $2, reminder_minute = $3, updated_at = now()`,
		chatID, hour, minute)
	return err
}

// parseClockTime accepts "19:00" or "9:00" and returns hour/minute in 0-23 / 0-59.
func parseClockTime(value string) (hour, minute int, err error) {
	parsed, err := time.Parse("15:04", strings.TrimSpace(value))
	if err != nil {
		return 0, 0, errors.New("используйте формат ЧЧ:ММ, например 19:00")
	}
	return parsed.Hour(), parsed.Minute(), nil
}

func (b *Bot) sendHomeworkList(ctx context.Context, chatID int64) error {
	rows, err := b.db.QueryContext(ctx, `SELECT id, subject, title, deadline FROM homework WHERE chat_id = $1 ORDER BY deadline ASC, id ASC`, chatID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var items []homework
	for rows.Next() {
		var item homework
		if err := rows.Scan(&item.ID, &item.Subject, &item.Title, &item.Deadline); err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(items) == 0 {
		return b.sendMessage(ctx, chatID, "Список заданий пуст.")
	}

	now := time.Now().In(b.config.TZ)
	var out strings.Builder
	out.WriteString("Задания по приоритету (ближайший дедлайн — выше):\n")
	for i, item := range items {
		deadline := item.Deadline.In(b.config.TZ)
		days := calendarDays(now, deadline, b.config.TZ)
		status := fmt.Sprintf("осталось %d дн.", days)
		if days < 0 {
			status = fmt.Sprintf("просрочено на %d дн.", -days)
		}
		if days == 0 {
			status = "срок сегодня"
		}
		fmt.Fprintf(&out, "\n%d. №%d — %s\n%s\nДедлайн: %s (%s)", i+1, item.ID, item.Subject, item.Title, deadline.Format("02.01.2006"), status)
	}
	return b.sendMessage(ctx, chatID, out.String())
}

func (b *Bot) reminderLoop(ctx context.Context) {
	// Tick once per minute, aligned to the start of the minute, and check which
	// chats have their configured reminder time (default 09:00) right now.
	now := time.Now()
	timer := time.NewTimer(time.Until(now.Truncate(time.Minute).Add(time.Minute)))
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			local := time.Now().In(b.config.TZ)
			b.sendDueReminders(ctx, local.Hour(), local.Minute())
			timer.Reset(time.Minute)
		}
	}
}

func (b *Bot) sendDueReminders(ctx context.Context, hour, minute int) {
	rows, err := b.db.QueryContext(ctx, `
		SELECT h.id, h.chat_id, h.subject, h.title, h.deadline
		FROM homework h
		LEFT JOIN user_settings s ON s.chat_id = h.chat_id
		WHERE COALESCE(s.reminder_hour, 9) = $1
		  AND COALESCE(s.reminder_minute, 0) = $2
		  AND h.deadline >= now() - interval '1 day'
		  AND h.deadline < now() + interval '8 days'
		ORDER BY h.deadline`, hour, minute)
	if err != nil {
		log.Printf("reminder query: %v", err)
		return
	}
	defer rows.Close()
	today := time.Now().In(b.config.TZ).Format("2006-01-02")
	for rows.Next() {
		var item homework
		var chatID int64
		if err := rows.Scan(&item.ID, &chatID, &item.Subject, &item.Title, &item.Deadline); err != nil {
			log.Print(err)
			continue
		}
		days := calendarDays(time.Now().In(b.config.TZ), item.Deadline.In(b.config.TZ), b.config.TZ)
		if days < 0 || days > 7 {
			continue
		}
		result, err := b.db.ExecContext(ctx, `INSERT INTO reminder_log (homework_id, reminder_date) VALUES ($1, $2) ON CONFLICT DO NOTHING`, item.ID, today)
		if err != nil {
			log.Print(err)
			continue
		}
		added, _ := result.RowsAffected()
		if added == 0 {
			continue
		}
		word := "дней"
		if days == 1 {
			word = "день"
		} else if days >= 2 && days <= 4 {
			word = "дня"
		}
		message := fmt.Sprintf("⏰ Напоминание\n%s — %s\nДо дедлайна осталось %d %s (%s).", item.Subject, item.Title, days, word, item.Deadline.In(b.config.TZ).Format("02.01.2006"))
		if days == 0 {
			message = fmt.Sprintf("⏰ Дедлайн сегодня!\n%s — %s", item.Subject, item.Title)
		}
		if err := b.sendMessage(ctx, chatID, message); err != nil {
			log.Printf("send reminder: %v", err)
		}
	}
}

func (b *Bot) sendMessage(ctx context.Context, chatID int64, text string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.config.Token)
	form := url.Values{"chat_id": {fmt.Sprint(chatID)}, "text": {text}}
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

func parseHomework(text string, tz *time.Location, now time.Time) (subject, title string, deadline time.Time, err error) {
	fields := make(map[string]string)
	for _, line := range strings.Split(text, "\n") {
		key, value, found := strings.Cut(line, ":")
		if found {
			fields[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
		}
	}
	subject, title = fields["предмет"], fields["задание"]
	dateText := fields["дедлайн"]
	if subject == "" || title == "" || dateText == "" {
		return "", "", time.Time{}, errors.New("нужны поля «Предмет», «Задание» и «Дедлайн»")
	}
	deadline, err = parseRussianDate(dateText, tz, now)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return subject, title, deadline, nil
}

func parseRussianDate(value string, tz *time.Location, now time.Time) (time.Time, error) {
	value = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(value, ".")))
	for _, layout := range []string{"02.01.2006", "02.01.06", "2006-01-02"} {
		if date, err := time.ParseInLocation(layout, value, tz); err == nil {
			return date, nil
		}
	}
	months := map[string]time.Month{"января": 1, "февраля": 2, "марта": 3, "апреля": 4, "мая": 5, "июня": 6, "июля": 7, "августа": 8, "сентября": 9, "октября": 10, "ноября": 11, "декабря": 12}
	parts := strings.Fields(value)
	if len(parts) < 2 || len(parts) > 3 {
		return time.Time{}, errors.New("используйте дату вида «6 октября» или «06.10.2026»")
	}
	var day, year int
	if _, err := fmt.Sscan(parts[0], &day); err != nil {
		return time.Time{}, errors.New("неверный день дедлайна")
	}
	month, ok := months[parts[1]]
	if !ok {
		return time.Time{}, errors.New("неверный месяц дедлайна")
	}
	year = now.Year()
	if len(parts) == 3 {
		if _, err := fmt.Sscan(parts[2], &year); err != nil {
			return time.Time{}, errors.New("неверный год дедлайна")
		}
	}
	date := time.Date(year, month, day, 0, 0, 0, 0, tz)
	if date.Month() != month || date.Day() != day {
		return time.Time{}, errors.New("такой даты не существует")
	}
	if len(parts) == 2 && date.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)) {
		date = date.AddDate(1, 0, 0)
	}
	return date, nil
}

func calendarDays(from, to time.Time, tz *time.Location) int {
	a := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, tz)
	b := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, tz)
	return int(b.Sub(a).Hours() / 24)
}

const formatExample = "Формат:\nПредмет: Разработка безопасного ПО\nЗадание: Практическая 1\nДедлайн: 6 октября"

const editExample = "Формат:\n/edit 3\nПредмет: Разработка безопасного ПО\nЗадание: Практическая 1\nДедлайн: 6 октября"

const helpText = "Присылайте задания в таком виде:\n\n" + formatExample + "\n\nКоманды:\n/list — задания по важности\n/edit НОМЕР — изменить задание (формат ниже)\n/delete НОМЕР — удалить выполненное задание\n/remindtime ЧЧ:ММ — во сколько напоминать каждый день (по умолчанию 09:00)\n\n" + editExample