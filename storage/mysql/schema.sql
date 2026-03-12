-- Parevo Core MySQL schema
-- Run this migration before using MySQL storage adapters

CREATE TABLE IF NOT EXISTS parevo_sessions (
    session_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    revoked TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_parevo_sessions_user_id ON parevo_sessions(user_id);
CREATE INDEX idx_parevo_sessions_revoked ON parevo_sessions(revoked);

CREATE TABLE IF NOT EXISTS parevo_refresh_tokens (
    token_id VARCHAR(255) PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    replaced_by VARCHAR(255),
    expires_at DATETIME(6) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_parevo_refresh_session ON parevo_refresh_tokens(session_id);
CREATE INDEX idx_parevo_refresh_user ON parevo_refresh_tokens(user_id);

-- Tenants: lifecycle
CREATE TABLE IF NOT EXISTS parevo_tenants (
    tenant_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    owner_id VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_parevo_tenants_status ON parevo_tenants(status);

-- Subject -> tenants mapping
CREATE TABLE IF NOT EXISTS parevo_subject_tenants (
    subject_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (subject_id, tenant_id)
);

CREATE INDEX idx_parevo_subject_tenants_subject ON parevo_subject_tenants(subject_id);

-- Permission grants
CREATE TABLE IF NOT EXISTS parevo_permission_grants (
    subject_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    permission VARCHAR(255) NOT NULL,
    PRIMARY KEY (subject_id, tenant_id, permission)
);

CREATE INDEX idx_parevo_permission_grants_subject_tenant ON parevo_permission_grants(subject_id, tenant_id);

-- Users: for social login
CREATE TABLE IF NOT EXISTS parevo_users (
    user_id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_parevo_users_email ON parevo_users(email);

-- Social accounts
CREATE TABLE IF NOT EXISTS parevo_social_accounts (
    provider VARCHAR(100) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    email_verified TINYINT(1) NOT NULL DEFAULT 0,
    name VARCHAR(255) NOT NULL DEFAULT '',
    avatar_url VARCHAR(512) NOT NULL DEFAULT '',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (provider, provider_user_id)
);

CREATE INDEX idx_parevo_social_accounts_user ON parevo_social_accounts(user_id);
