CREATE TABLE IF NOT EXISTS daily_recommendation_log (
    chat_id BIGINT NOT NULL,
    recommendation_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'sent')),
    attempts SMALLINT NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_until TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    last_error TEXT,
    PRIMARY KEY (chat_id, recommendation_date)
);

CREATE INDEX IF NOT EXISTS daily_recommendation_pending_idx
    ON daily_recommendation_log (next_attempt_at)
    WHERE status <> 'sent';
