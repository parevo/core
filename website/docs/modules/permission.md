# Permission Module

Permission check service with wildcard support.

## Usage

```go
permStore := &memory.PermissionStore{
    Grants: map[string]bool{
        "user-1|tenant-a|orders:read": true,
    },
}
svc := permission.NewService(permStore)
err := svc.Check(ctx, subject, tenantID, "orders:read")
```
