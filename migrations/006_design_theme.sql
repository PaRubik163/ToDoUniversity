ALTER TABLE user_settings
    ADD COLUMN IF NOT EXISTS design_theme TEXT NOT NULL DEFAULT 'classic'
    CHECK (design_theme IN ('classic', 'neon'));
