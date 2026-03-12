# Auth Module

JWT-based authentication, guards, middleware adapters.

## Submodules

- `auth/mfa` — TOTP 2FA, recovery codes
- `auth/apikey` — API key validation
- `auth/webauthn` — WebAuthn/Passkeys
- `auth/magiclink` — magic link / email OTP
- `auth/blacklist` — JWT blacklist
- `auth/oauth2provider` — OAuth2 authorization server

## Config

- `Issuer`, `Audience`, `SecretKey` (or `SigningKeys` + `ActiveKID`)
- `AccessTokenTTL`, `RefreshTokenTTL`
