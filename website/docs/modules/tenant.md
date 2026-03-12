# Tenant Module

Tenant selection, override policy, lifecycle (create/suspend/delete).

## Usage

```go
tenantStore := &memory.TenantStore{
    SubjectTenants: map[string][]string{
        "user-1": {"tenant-a", "tenant-b"},
    },
}
svc := tenant.NewService(tenantStore)
tenantID, err := svc.Resolve(ctx, subject, defaultTenant, requestedTenant, policy)
```
