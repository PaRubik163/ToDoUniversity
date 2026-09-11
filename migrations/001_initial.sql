CREATE TABLE IF NOT EXISTS homework (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    subject TEXT NOT NULL,
    title TEXT NOT NULL,
    deadline TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS homework_chat_deadline_idx
    ON homework (chat_id, deadline);

CREATE TABLE IF NOT EXISTS reminder_log (
    homework_id BIGINT NOT NULL REFERENCES homework(id) ON DELETE CASCADE,
    reminder_date DATE NOT NULL,
    PRIMARY KEY (homework_id, reminder_date)
);
