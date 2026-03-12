# Folder Structure

This repo follows a public library layout; packages live at the root by design.

- `auth/`: token, middleware, guard, audit, ratelimit
- `tenant/`: tenant selection and override rules
- `permission/`: permission check services
- `social/`: social callback + account linking
- `storage/`: DB adapter contracts and in-memory implementations
- `notification/`: email, SMS, WebSocket sender interface and adapters
- `blob/`: object storage (S3, R2) interface and adapters
- `examples/`: working integration examples

Each concern is a separate package. A `pkg/` subfolder is not used since the library is imported directly.
