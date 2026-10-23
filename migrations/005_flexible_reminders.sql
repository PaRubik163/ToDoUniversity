ALTER TABLE user_settings
    ADD COLUMN IF NOT EXISTS reminder_days SMALLINT[] NOT NULL DEFAULT ARRAY[7]::SMALLINT[],
    ADD COLUMN IF NOT EXISTS reminder_times TIME[] NOT NULL DEFAULT ARRAY['09:00'::TIME];

UPDATE user_settings
SET reminder_times = ARRAY[make_time(reminder_hour, reminder_minute, 0)];

ALTER TABLE reminder_log
    ADD COLUMN IF NOT EXISTS reminder_time TIME NOT NULL DEFAULT '09:00';

ALTER TABLE reminder_log
    DROP CONSTRAINT IF EXISTS reminder_log_pkey;

ALTER TABLE reminder_log
    ADD CONSTRAINT reminder_log_pkey PRIMARY KEY (homework_id, reminder_date, reminder_time);
