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

-- Tenants: lifecycle (Create, Suspend, Resume, Delete, Status, ListTenants)
CREATE TABLE IF NOT EXISTS parevo_tenants (
    tenant_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    owner_id VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_parevo_tenants_status ON parevo_tenants(status);

-- Tenant: subject -> tenants mapping (for ResolveSubjectTenants)
CREATE TABLE IF NOT EXISTS parevo_subject_tenants (
    subject_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (subject_id, tenant_id)
);
CREATE INDEX IF NOT EXISTS idx_parevo_subject_tenants_subject ON parevo_subject_tenants(subject_id);

-- Permissions: subject + tenant -> permission grants
CREATE TABLE IF NOT EXISTS parevo_permission_grants (
    subject_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    permission VARCHAR(255) NOT NULL,
    PRIMARY KEY (subject_id, tenant_id, permission)
);
CREATE INDEX IF NOT EXISTS idx_parevo_permission_grants_subject_tenant ON parevo_permission_grants(subject_id, tenant_id);

-- Users: for social login (FindOrCreateUserByEmail)
CREATE TABLE IF NOT EXISTS parevo_users (
    user_id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_parevo_users_email ON parevo_users(email);

-- Social accounts: provider + provider_user_id -> user_id
CREATE TABLE IF NOT EXISTS parevo_social_accounts (
    provider VARCHAR(100) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL DEFAULT '',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL DEFAULT '',
    avatar_url VARCHAR(512) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider, provider_user_id)
);
CREATE INDEX IF NOT EXISTS idx_parevo_social_accounts_user ON parevo_social_accounts(user_id);
