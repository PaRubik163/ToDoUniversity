CREATE TABLE IF NOT EXISTS homework_status (
    homework_id BIGINT PRIMARY KEY REFERENCES homework(id) ON DELETE CASCADE,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS homework_status_completed_at_idx
    ON homework_status (completed_at)
    WHERE is_completed = TRUE;

INSERT INTO homework_status (homework_id, is_completed)
SELECT id, FALSE FROM homework
ON CONFLICT (homework_id) DO NOTHING;
