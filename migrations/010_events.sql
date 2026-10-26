CREATE TABLE IF NOT EXISTS calendar_events (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    event_at TIMESTAMPTZ NOT NULL,
    reminder_hours SMALLINT NOT NULL DEFAULT 1 CHECK (reminder_hours BETWEEN 0 AND 168),
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS calendar_events_user_date_idx
    ON calendar_events (chat_id, event_at);

CREATE TABLE IF NOT EXISTS calendar_event_reminder_log (
    event_id BIGINT PRIMARY KEY REFERENCES calendar_events(id) ON DELETE CASCADE,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
