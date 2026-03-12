# Product Gaps / Roadmap

The library provides a solid foundation. Priority gaps for production:

## P0 - Security and Compliance

- OIDC/Social Login callback + Google provider adapter (done)
- Key rotation (KID-based JWT) (done)
- Refresh token rotation + reuse detection (done)
- Session revoke store, revoke check, logout-all-by-user (done)
- Rate limiting + brute-force lockout manager (done)
- Audit log hook + structured logger integration (done)

## P1 - Multi-tenant Runtime

- Tenant lifecycle API (done)
- Tenant-level feature flags and plan/policy binding (done)
- Permission wildcard support (done)
- ABAC conditions (resource owner, department, environment) (done)

## P2 - DX and Operations

- Config validation package extensions (strict mode) (done)
- OpenTelemetry tracing hooks (done)
- Prometheus metrics (done)
- Migration-ready SQL adapter packages - postgres (done)
- Versioned changelog + release automation (done)

## P2 - Security & Auth Extensions (done)

- WebAuthn/Passkeys
- Magic link / Email OTP
- MFA recovery codes
- OAuth2 Provider mode
- OAuth2 scopes in claims

## P3 - Operations (done)

- Session metadata (IP, user-agent, last activity)
- Audit log search/export
- IP allowlist/blocklist
- JWT blacklist
- Consent management

## P4 - Enterprise & Integrations (done)

- SAML 2.0
- LDAP/Active Directory
- Webhook events
- GraphQL adapter
