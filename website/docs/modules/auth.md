# Auth Module

JWT-based authentication, guards, middleware adapters.

## Submodules

- `auth/mfa` — TOTP 2FA, recovery codes
- `auth/apikey` — API key validation (`auth/apikey/memory` for in-memory store)
- `auth/webauthn` — WebAuthn/Passkeys
- `auth/magiclink` — magic link / email OTP
- `auth/blacklist` — JWT blacklist (logout). Wire via `auth.Modules.Blacklist`; `ParseAndValidate` rejects blacklisted tokens.
- `auth/oauth2provider` — OAuth2 authorization server

## Config

- `Issuer`, `Audience`, `SecretKey` (or `SigningKeys` + `ActiveKID`)
- `AccessTokenTTL`, `RefreshTokenTTL`
