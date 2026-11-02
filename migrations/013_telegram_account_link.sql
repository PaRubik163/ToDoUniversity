ALTER TABLE auth_refresh_tokens
    DROP CONSTRAINT IF EXISTS auth_refresh_tokens_user_id_fkey;

ALTER TABLE auth_refresh_tokens
    ADD CONSTRAINT auth_refresh_tokens_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES auth_accounts(id)
    ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS auth_link_codes (
    token_hash BYTEA PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES auth_accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
