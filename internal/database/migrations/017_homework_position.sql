ALTER TABLE homework ADD COLUMN IF NOT EXISTS position BIGINT NOT NULL DEFAULT 0;

WITH ranked AS (
    SELECT id, row_number() OVER (PARTITION BY chat_id ORDER BY deadline ASC, id ASC) - 1 AS position
    FROM homework
)
UPDATE homework h SET position = ranked.position
FROM ranked WHERE ranked.id = h.id AND h.position = 0;

CREATE INDEX IF NOT EXISTS homework_chat_position_idx ON homework (chat_id, position, id);
