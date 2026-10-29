ALTER TABLE homework
    ADD COLUMN IF NOT EXISTS priority TEXT NOT NULL DEFAULT 'medium',
    ADD COLUMN IF NOT EXISTS note TEXT NOT NULL DEFAULT '';

ALTER TABLE homework
    ADD CONSTRAINT homework_priority_check
    CHECK (priority IN ('low', 'medium', 'high'));
