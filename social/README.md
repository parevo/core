# social package

The `social` package normalizes provider callback flows:

- Provider code exchange
- Social account linking
- Access token issuance

This layer does not depend on provider SDKs directly; it expects provider-specific adapters.

Example providers: `social/providers/google`, `social/providers/github`.
