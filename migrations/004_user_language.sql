ALTER TABLE user_settings
    ADD COLUMN IF NOT EXISTS language TEXT NOT NULL DEFAULT 'ru'
    CHECK (language IN ('ru', 'en'));
