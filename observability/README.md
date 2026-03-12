# observability

This folder contains production-grade observability components.

## logging

`observability/logging` provides a structured logger:

- Development mode: readable key=value output
- Production mode: JSON output
- Minimum log level
- Sensitive key redaction (`token`, `password`, `secret`, etc.)
- Context enrichment from `request_id`, `tenant_id`, `user_id`
