package tenantsql

import (
	"fmt"
	"strings"

	"github.com/parevo/core/auth"
)

func EnsureTenant(tenantID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return auth.ErrMissingTenant
	}
	return nil
}

func AppendTenantFilter(baseQuery string, tenantID string, args ...any) (string, []any, error) {
	if err := EnsureTenant(tenantID); err != nil {
		return "", nil, err
	}

	query := strings.TrimSpace(baseQuery)
	if query == "" {
		return "", nil, fmt.Errorf("%w: empty query", auth.ErrInvalidConfig)
	}

	lower := strings.ToLower(query)
	var filtered string
	if strings.Contains(lower, " where ") {
		filtered = query + " AND tenant_id = ?"
	} else {
		filtered = query + " WHERE tenant_id = ?"
	}
	newArgs := append([]any{}, args...)
	newArgs = append(newArgs, tenantID)
	return filtered, newArgs, nil
}
