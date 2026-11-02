ALTER TABLE user_settings
    DROP CONSTRAINT IF EXISTS user_settings_language_check;

ALTER TABLE user_settings
    ADD CONSTRAINT user_settings_language_check
    CHECK (language IN ('ru', 'en', 'zh', 'es', 'ar', 'hy'));
