CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              token_hash VARCHAR(64) PRIMARY KEY,
                                              user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                                              expires_at TIMESTAMP NOT NULL,
                                              created_at TIMESTAMP NOT NULL
);