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
- Tenant-level feature flags and plan/policy binding
- Permission wildcard support (done)
- ABAC conditions (resource owner, department, environment)

## P2 - DX and Operations

- Config validation package extensions (strict mode)
- OpenTelemetry tracing hooks (done)
- Prometheus metrics (done)
- Migration-ready SQL adapter packages - postgres (done)
- Versioned changelog + release automation (done)
