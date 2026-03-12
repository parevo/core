package mysql

import (
	"context"
	"database/sql"

	"github.com/parevo/core/storage"
)

// TenantStore implements storage.TenantStore for MySQL.
type TenantStore struct {
	db *sql.DB
}

// NewTenantStore creates a MySQL TenantStore.
func NewTenantStore(db *sql.DB) *TenantStore {
	return &TenantStore{db: db}
}

// ResolveSubjectTenants returns tenant IDs the subject has access to.
func (s *TenantStore) ResolveSubjectTenants(ctx context.Context, subjectID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT tenant_id FROM parevo_subject_tenants WHERE subject_id = ?`,
		subjectID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tenants []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, rows.Err()
}

var _ storage.TenantStore = (*TenantStore)(nil)

// TenantLifecycleStore implements storage.TenantLifecycleStore and storage.TenantListStore for MySQL.
type TenantLifecycleStore struct {
	db *sql.DB
}

// NewTenantLifecycleStore creates a MySQL TenantLifecycleStore.
func NewTenantLifecycleStore(db *sql.DB) *TenantLifecycleStore {
	return &TenantLifecycleStore{db: db}
}

// Create creates a new tenant.
func (s *TenantLifecycleStore) Create(ctx context.Context, tenantID, name, ownerID string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT IGNORE INTO parevo_tenants (tenant_id, name, owner_id, status) VALUES (?, ?, ?, 'active')`,
		tenantID, name, ownerID)
	return err
}

// Suspend marks the tenant as suspended.
func (s *TenantLifecycleStore) Suspend(ctx context.Context, tenantID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_tenants SET status = 'suspended' WHERE tenant_id = ?`, tenantID)
	return err
}

// Resume marks the tenant as active.
func (s *TenantLifecycleStore) Resume(ctx context.Context, tenantID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_tenants SET status = 'active' WHERE tenant_id = ?`, tenantID)
	return err
}

// Delete marks the tenant as deleted (soft delete).
func (s *TenantLifecycleStore) Delete(ctx context.Context, tenantID string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE parevo_tenants SET status = 'deleted' WHERE tenant_id = ?`, tenantID)
	return err
}

// Status returns the tenant status.
func (s *TenantLifecycleStore) Status(ctx context.Context, tenantID string) (storage.TenantStatus, error) {
	var status string
	err := s.db.QueryRowContext(ctx, `SELECT status FROM parevo_tenants WHERE tenant_id = ?`, tenantID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return storage.TenantStatus(status), nil
}

// ListTenants returns all non-deleted tenants.
func (s *TenantLifecycleStore) ListTenants(ctx context.Context) ([]storage.TenantInfo, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT tenant_id, status FROM parevo_tenants WHERE status != 'deleted'`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []storage.TenantInfo
	for rows.Next() {
		var info storage.TenantInfo
		if err := rows.Scan(&info.ID, &info.Status); err != nil {
			return nil, err
		}
		out = append(out, info)
	}
	return out, rows.Err()
}

var _ storage.TenantLifecycleStore = (*TenantLifecycleStore)(nil)
var _ storage.TenantListStore = (*TenantLifecycleStore)(nil)
