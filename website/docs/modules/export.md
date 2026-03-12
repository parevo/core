# Export Module

GDPR/compliance veri dışa aktarma (data portability).

## Usage

```go
import "github.com/parevo/core/export"

payload := export.NewPayload(userID)
payload.Profile = map[string]any{"email": "u@example.com", "name": "User"}
payload.Sessions = []map[string]any{{"id": "s1", "created_at": "..."}}
payload.Consents = []map[string]any{{"client_id": "c1", "scopes": []string{"openid"}}}

// JSON
jsonBytes, _ := export.ToJSON(payload)

// Blob storage
export.ToBlob(ctx, blobStore, "exports", "user-123.json", payload)
```

## Payload Alanları

- `Profile` — kullanıcı profil verisi
- `Sessions` — oturum listesi
- `Consents` — OAuth consent kayıtları
- `Permissions` — izin listesi

Uygulama kendi storage'ından veriyi çekip payload'a doldurur.
