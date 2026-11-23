ALTER TABLE homework
    ADD COLUMN IF NOT EXISTS recurrence TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS recurrence_days SMALLINT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS recurrence_source_id BIGINT REFERENCES homework(id) ON DELETE SET NULL;

ALTER TABLE homework
    ADD CONSTRAINT homework_recurrence_check
    CHECK (recurrence IN ('', 'daily', 'weekly', 'monthly'));

ALTER TABLE homework
    ADD CONSTRAINT homework_recurrence_days_check
    CHECK (recurrence_days <@ ARRAY[1,2,3,4,5,6,7]::smallint[]);

CREATE UNIQUE INDEX IF NOT EXISTS homework_recurrence_source_idx
    ON homework (recurrence_source_id) WHERE recurrence_source_id IS NOT NULL;
