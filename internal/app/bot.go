package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"homework-reminder-bot/internal/config"
	"homework-reminder-bot/internal/domain"
	"homework-reminder-bot/internal/ratelimit"
	"homework-reminder-bot/internal/repository/postgres"
	"homework-reminder-bot/internal/service"
)

const pollTimeout = 25 * time.Second

type Bot struct {
	config  config.Config
	db      *sql.DB
	client  *http.Client
	limiter ratelimit.Limiter
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
		ID   int64  `json:"id"`
		Type string `json:"type"`
	} `json:"chat"`
	From struct {
		ID           int64  `json:"id"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Text string `json:"text"`
}

type homework = domain.Homework

// New constructs the Telegram bot and its HTTP application services.
func New(appConfig config.Config, db *sql.DB, limiter ratelimit.Limiter) *Bot {
	return &Bot{
		config:  appConfig,
		db:      db,
		client:  &http.Client{Timeout: pollTimeout + 10*time.Second},
		limiter: limiter,
	}
}

// Run starts background reminders, the Mini App server, and Telegram polling.
func (b *Bot) Run(ctx context.Context) error {
	if b.config.MiniAppURL != "" {
		if err := b.setMenuButton(ctx); err != nil {
			log.Printf("set menu button: %v", err)
		}
	} else {
		log.Println("MINI_APP_URL not set: skipping Telegram menu button, mini app will only be reachable directly by URL")
	}

	go b.reminderLoop(ctx)
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- b.startWebServer(ctx, ":"+b.config.Port) }()

	pollErrors := make(chan error, 1)
	go func() { pollErrors <- b.poll(ctx) }()

	log.Println("bot is running")
	select {
	case <-ctx.Done():
		return nil
	case err := <-serverErrors:
		if ctx.Err() != nil || err == nil {
			return nil
		}
		return fmt.Errorf("web server: %w", err)
	case err := <-pollErrors:
		if ctx.Err() != nil || err == nil {
			return nil
		}
		return fmt.Errorf("telegram polling: %w", err)
	}
}

func (b *Bot) poll(ctx context.Context) error {
	var offset int64
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		updates, err := b.getUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("getUpdates: %v", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
			}
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
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", b.config.TelegramBotToken)
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
	case text == "/start":
		if _, err := b.db.ExecContext(ctx, `INSERT INTO user_settings (chat_id, language) VALUES ($1, $2) ON CONFLICT (chat_id) DO NOTHING`, msg.Chat.ID, normalizedLanguage(msg.From.LanguageCode)); err != nil {
			log.Printf("save Telegram language: %v", err)
		}
		welcome, button := startCopy(msg.From.LanguageCode)
		if b.config.MiniAppURL != "" {
			if err := b.sendWebAppButton(ctx, msg.Chat.ID, welcome, button, b.config.MiniAppURL); err == nil {
				return nil
			} else {
				log.Printf("send web app button: %v", err)
			}
		}
		return b.sendMessage(ctx, msg.Chat.ID, welcome)
	case text == "/help":
		if b.config.MiniAppURL != "" {
			_, button := startCopy(msg.From.LanguageCode)
			err := b.sendWebAppButton(ctx, msg.Chat.ID, helpText, button, b.config.MiniAppURL)
			if err == nil {
				return nil
			}
			log.Printf("send web app button: %v", err)
		}
		return b.sendMessage(ctx, msg.Chat.ID, helpText)
	case strings.HasPrefix(text, "/link"):
		parts := strings.Fields(text)
		if msg.Chat.Type != "private" || msg.From.ID <= 0 || len(parts) != 2 {
			return b.sendMessage(ctx, msg.Chat.ID, "Для привязки отправьте боту в личном чате: /link КОД")
		}
		allowed, retryAfter, err := b.limiter.Allow(ctx, "telegram:link:user:"+strconv.FormatInt(msg.From.ID, 10), 5, 10*time.Minute)
		if err != nil {
			log.Printf("Telegram link rate limit unavailable: %v", err)
			return b.sendMessage(ctx, msg.Chat.ID, "Сервис временно недоступен. Попробуйте позже.")
		}
		if !allowed {
			seconds := int(retryAfter.Round(time.Second).Seconds())
			return b.sendMessage(ctx, msg.Chat.ID, fmt.Sprintf("Слишком много попыток. Попробуйте через %d секунд.", seconds))
		}
		linker := service.NewAuthService(postgres.NewAuthRepository(b.db), b.config.JWTSecret)
		if err := linker.LinkTelegramAccount(ctx, parts[1], msg.From.ID); err != nil {
			return b.sendMessage(ctx, msg.Chat.ID, "Код привязки недействителен или истёк. Создайте новый на сайте.")
		}
		return b.sendMessage(ctx, msg.Chat.ID, "Аккаунты связаны. Обновите страницу сайта — там появятся ваши прежние задачи Telegram.")
	case text == "/list" || text == "Список":
		return b.sendHomeworkList(ctx, msg.Chat.ID)
	case strings.HasPrefix(text, "/delete") || strings.HasPrefix(text, "Удалить"):
		return b.deleteHomework(ctx, msg.Chat.ID, text)
	case strings.HasPrefix(text, "/edit"):
		return b.editHomeworkText(ctx, msg.Chat.ID, text)
	case strings.HasPrefix(text, "/remindtime"):
		return b.setReminderTimeText(ctx, msg.Chat.ID, text)
	case strings.Contains(strings.ToLower(text), "предмет:"):
		subject, title, deadline, err := parseHomework(text, b.config.Timezone, time.Now().In(b.config.Timezone))
		if err != nil {
			return b.sendMessage(ctx, msg.Chat.ID, "Не удалось прочитать задание: "+err.Error()+"\n\n"+formatExample)
		}
		var id int64
		err = b.db.QueryRowContext(ctx, `
			WITH new_homework AS (
				INSERT INTO homework (chat_id, subject, title, deadline)
				VALUES ($1, $2, $3, $4)
				RETURNING id
			)
			INSERT INTO homework_status (homework_id, is_completed)
			SELECT id, FALSE FROM new_homework
			RETURNING homework_id`, msg.Chat.ID, subject, title, deadline).Scan(&id)
		if err != nil {
			return err
		}
		return b.sendMessage(ctx, msg.Chat.ID, fmt.Sprintf("Записала задание №%d. Напомню за 7 дней до дедлайна и далее каждый день.", id))
	default:
		return b.sendMessage(ctx, msg.Chat.ID, "Я понимаю добавление задания, /list, /edit, /delete и /remindtime.\n\n"+formatExample)
	}
}

// startCopy uses Telegram's language_code from the /start update. initData is
// only available after the user has opened the Mini App.
func startCopy(languageCode string) (message, button string) {
	language := normalizedLanguage(languageCode)
	switch language {
	case "ru":
		return "Добавьте задачу и выберите время напоминания — бот поможет ничего не забыть. Чтобы начать, нажмите кнопку ниже.", "Открыть приложение"
	case "en":
		return "Add a task, choose when you’d like to be reminded, and the bot will help you stay on top of it. Tap the button below to get started.", "Open app"
	case "zh":
		return "添加任务并选择提醒时间，机器人会帮你及时记住每件事。点击下方按钮开始使用。", "打开应用"
	case "es":
		return "Añade una tarea y elige cuándo quieres recibir un recordatorio. El bot te ayudará a tenerlo todo bajo control. Pulsa el botón para empezar.", "Abrir la app"
	case "ar":
		return "أضف مهمة واختر موعد التذكير، وسيساعدك البوت على متابعة كل ما عليك إنجازه. اضغط الزر أدناه للبدء.", "فتح التطبيق"
	case "hy":
		return "Ավելացրեք առաջադրանք և ընտրեք հիշեցման ժամը․ բոտը կօգնի ձեզ ամեն ինչ ժամանակին հիշել։ Սկսելու համար սեղմեք ներքևի կոճակը։", "Բացել հավելվածը"
	default:
		return "Add a task, choose when you’d like to be reminded, and the bot will help you stay on top of it. Tap the button below to get started.", "Open app"
	}
}

func normalizedLanguage(languageCode string) string {
	language := strings.ToLower(strings.TrimSpace(languageCode))
	if i := strings.IndexAny(language, "-_"); i >= 0 {
		language = language[:i]
	}
	if validLanguage(language) {
		return language
	}
	return "en"
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
	subject, title, deadline, err := parseHomework(lines[1], b.config.Timezone, time.Now().In(b.config.Timezone))
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
		INSERT INTO user_settings (chat_id, reminder_hour, reminder_minute, reminder_times, updated_at)
		VALUES ($1, $2, $3, ARRAY[make_time($2, $3, 0)], now())
		ON CONFLICT (chat_id) DO UPDATE SET reminder_hour = $2, reminder_minute = $3,
			reminder_times = ARRAY[make_time($2, $3, 0)], updated_at = now()`,
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
	rows, err := b.db.QueryContext(ctx, `
		SELECT h.id, h.subject, h.title, h.deadline
		FROM homework h
		LEFT JOIN homework_status s ON s.homework_id = h.id
		WHERE h.chat_id = $1 AND COALESCE(s.is_completed, FALSE) = FALSE
		ORDER BY h.deadline ASC, h.id ASC`, chatID)
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

	now := time.Now().In(b.config.Timezone)
	var out strings.Builder
	out.WriteString("Задания по приоритету (ближайший дедлайн — выше):\n")
	for i, item := range items {
		deadline := item.Deadline.In(b.config.Timezone)
		days := domain.CalendarDays(now, deadline, b.config.Timezone)
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
			local := time.Now().In(b.config.Timezone)
			b.sendDueReminders(ctx, local.Hour(), local.Minute())
			b.sendDueEventReminders(ctx, local)
			timer.Reset(time.Minute)
		}
	}
}

func (b *Bot) sendDueEventReminders(ctx context.Context, now time.Time) {
	rows, err := b.db.QueryContext(ctx, `
		SELECT id, chat_id, title, event_at
		FROM calendar_events
		WHERE is_completed = FALSE AND reminder_hours > 0
		  AND event_at - reminder_hours * interval '1 hour' <= $1
		  AND event_at > $1
		  AND NOT EXISTS (
			  SELECT 1 FROM calendar_event_reminder_log sent WHERE sent.event_id = calendar_events.id
		  )`, now)
	if err != nil {
		log.Printf("event reminder query: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id, chatID int64
		var title string
		var eventAt time.Time
		if err := rows.Scan(&id, &chatID, &title, &eventAt); err != nil {
			log.Printf("scan event reminder: %v", err)
			continue
		}
		message := fmt.Sprintf("📅 Напоминание о событии\n%s\n%s", title, eventAt.In(b.config.Timezone).Format("02.01.2006 в 15:04"))
		if err := b.sendMessage(ctx, chatID, message); err != nil {
			log.Printf("send event reminder: %v", err)
			continue
		}
		if _, err := b.db.ExecContext(ctx, `INSERT INTO calendar_event_reminder_log (event_id) VALUES ($1) ON CONFLICT DO NOTHING`, id); err != nil {
			log.Printf("log event reminder: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("iterate event reminders: %v", err)
	}
}

func (b *Bot) sendDueReminders(ctx context.Context, hour, minute int) {
	rows, err := b.db.QueryContext(ctx, `
		SELECT h.id, h.chat_id, h.subject, h.title, h.deadline,
			array_to_json(COALESCE(h.reminder_days, ARRAY[7]::smallint[]))::text
		FROM homework h
		LEFT JOIN user_settings s ON s.chat_id = h.chat_id
		LEFT JOIN homework_status hs ON hs.homework_id = h.id
		WHERE EXISTS (
			SELECT 1 FROM unnest(COALESCE(s.reminder_times, ARRAY['09:00'::TIME])) AS configured(time_value)
			WHERE configured.time_value = make_time($1, $2, 0)
		  )
		  AND COALESCE(hs.is_completed, FALSE) = FALSE
		  AND h.deadline < now() + interval '21 days'
		ORDER BY h.deadline`, hour, minute)
	if err != nil {
		log.Printf("reminder query: %v", err)
		return
	}
	defer rows.Close()
	today := time.Now().In(b.config.Timezone).Format("2006-01-02")
	reminderTime := fmt.Sprintf("%02d:%02d", hour, minute)
	for rows.Next() {
		var item homework
		var chatID int64
		var reminderDaysJSON string
		if err := rows.Scan(&item.ID, &chatID, &item.Subject, &item.Title, &item.Deadline, &reminderDaysJSON); err != nil {
			log.Print(err)
			continue
		}
		days := domain.CalendarDays(time.Now().In(b.config.Timezone), item.Deadline.In(b.config.Timezone), b.config.Timezone)
		if days > 0 {
			var configuredDays []int
			if err := json.Unmarshal([]byte(reminderDaysJSON), &configuredDays); err != nil {
				log.Printf("decode reminder days: %v", err)
				continue
			}
			selected := false
			for _, configuredDay := range configuredDays {
				if days == configuredDay {
					selected = true
					break
				}
			}
			if !selected {
				continue
			}
		}
		result, err := b.db.ExecContext(ctx, `INSERT INTO reminder_log (homework_id, reminder_date, reminder_time) VALUES ($1, $2, $3::time) ON CONFLICT DO NOTHING`, item.ID, today, reminderTime)
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
		message := fmt.Sprintf("⏰ Напоминание\n%s — %s\nДо дедлайна осталось %d %s (%s).", item.Subject, item.Title, days, word, item.Deadline.In(b.config.Timezone).Format("02.01.2006"))
		if days < 0 {
			overdueDays := -days
			overdueWord := "дней"
			if overdueDays == 1 {
				overdueWord = "день"
			} else if overdueDays >= 2 && overdueDays <= 4 {
				overdueWord = "дня"
			}
			message = fmt.Sprintf("⚠️ Дедлайн истек %d %s назад\n%s — %s (дедлайн: %s).", overdueDays, overdueWord, item.Subject, item.Title, item.Deadline.In(b.config.Timezone).Format("02.01.2006"))
		} else if days == 0 {
			message = fmt.Sprintf("⏰ Дедлайн сегодня!\n%s — %s", item.Subject, item.Title)
		}
		if err := b.sendMessage(ctx, chatID, message); err != nil {
			log.Printf("send reminder: %v", err)
		}
	}
}

func (b *Bot) sendMessage(ctx context.Context, chatID int64, text string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.config.TelegramBotToken)
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

const formatExample = "Формат:\nПредмет: Разработка безопасного ПО\nЗадание: Практическая 1\nДедлайн: 6 октября"

const editExample = "Формат:\n/edit 3\nПредмет: Разработка безопасного ПО\nЗадание: Практическая 1\nДедлайн: 6 октября"

const helpText = "Присылайте задания в таком виде:\n\n" + formatExample + "\n\nКоманды:\n/list — задания по важности\n/edit НОМЕР — изменить задание (формат ниже)\n/delete НОМЕР — удалить выполненное задание\n/remindtime ЧЧ:ММ — во сколько напоминать каждый день (по умолчанию 09:00)\n\n" + editExample
