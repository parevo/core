package main

import (
	"context"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/permission"
	"github.com/parevo/core/storage"
	"github.com/parevo/core/tenant"
)

// Example adapter: map your own DB schema to Parevo Core storage interfaces.
type tenantRepo struct{}

func (r *tenantRepo) ResolveSubjectTenants(_ context.Context, subjectID string) ([]string, error) {
	// Query your own tables/collections here.
	if subjectID == "user-1" {
		return []string{"tenant-a"}, nil
	}
	return []string{}, nil
}

type permissionRepo struct{}

func (r *permissionRepo) HasPermission(_ context.Context, subjectID, tenantID, perm string, _ []string) (bool, error) {
	// Query your own permission model here.
	return subjectID == "user-1" && tenantID == "tenant-a" && perm == "orders:read", nil
}

func main() {
	tenantSvc := tenant.NewService(&tenantRepo{})
	permSvc := permission.NewService(&permissionRepo{})

	_, _ = auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		Tenant:     tenantSvc,
		Permission: permSvc,
	})

	_ = storage.Subject{}
}
