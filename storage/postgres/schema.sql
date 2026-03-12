-- Parevo Core postgres schema
-- Run this migration before using PostgresSessionStore or PostgresRefreshStore

CREATE TABLE IF NOT EXISTS parevo_sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_parevo_sessions_user_id ON parevo_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_parevo_sessions_revoked ON parevo_sessions(revoked) WHERE revoked = TRUE;

CREATE TABLE IF NOT EXISTS parevo_refresh_tokens (
    token_id VARCHAR(255) PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    replaced_by VARCHAR(255),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_parevo_refresh_session ON parevo_refresh_tokens(session_id);
CREATE INDEX IF NOT EXISTS idx_parevo_refresh_user ON parevo_refresh_tokens(user_id);
