CREATE TABLE IF NOT EXISTS user_settings (
    chat_id BIGINT PRIMARY KEY,
    reminder_hour SMALLINT NOT NULL DEFAULT 9 CHECK (reminder_hour BETWEEN 0 AND 23),
    reminder_minute SMALLINT NOT NULL DEFAULT 0 CHECK (reminder_minute BETWEEN 0 AND 59),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);