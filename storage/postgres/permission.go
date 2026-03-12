package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/parevo/core/storage"
)

// PermissionStore implements storage.PermissionStore and storage.PermissionGrantStore for Postgres.
type PermissionStore struct {
	db *sql.DB
}

// NewPermissionStore creates a Postgres PermissionStore.
func NewPermissionStore(db *sql.DB) *PermissionStore {
	return &PermissionStore{db: db}
}

// HasPermission checks if the subject has the permission in the tenant.
func (s *PermissionStore) HasPermission(ctx context.Context, subjectID, tenantID, permission string, _ []string) (bool, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT permission FROM parevo_permission_grants WHERE subject_id = $1 AND tenant_id = $2`,
		subjectID, tenantID)
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var granted string
		if err := rows.Scan(&granted); err != nil {
			return false, err
		}
		if matchPermission(granted, permission) {
			return true, nil
		}
	}
	return false, rows.Err()
}

func matchPermission(granted, requested string) bool {
	if granted == requested || granted == "*" || granted == "*:*" {
		return true
	}
	gParts := strings.Split(granted, ":")
	rParts := strings.Split(requested, ":")
	if len(gParts) != 2 || len(rParts) != 2 {
		return false
	}
	if gParts[0] == "*" && gParts[1] == "*" {
		return true
	}
	if gParts[0] == rParts[0] && gParts[1] == "*" {
		return true
	}
	if gParts[0] == "*" && gParts[1] == rParts[1] {
		return true
	}
	return false
}

// Grant adds a permission grant.
func (s *PermissionStore) Grant(ctx context.Context, subjectID, tenantID, permission string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parevo_permission_grants (subject_id, tenant_id, permission) VALUES ($1, $2, $3) ON CONFLICT (subject_id, tenant_id, permission) DO NOTHING`,
		subjectID, tenantID, permission)
	return err
}

// Revoke removes a permission grant.
func (s *PermissionStore) Revoke(ctx context.Context, subjectID, tenantID, permission string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM parevo_permission_grants WHERE subject_id = $1 AND tenant_id = $2 AND permission = $3`,
		subjectID, tenantID, permission)
	return err
}

// ListGrants returns all permissions for the subject in the tenant.
func (s *PermissionStore) ListGrants(ctx context.Context, subjectID, tenantID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT permission FROM parevo_permission_grants WHERE subject_id = $1 AND tenant_id = $2`,
		subjectID, tenantID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

var _ storage.PermissionStore = (*PermissionStore)(nil)
var _ storage.PermissionGrantStore = (*PermissionStore)(nil)
