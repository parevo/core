package tenantsql

import (
	"errors"
	"testing"

	"github.com/parevo/core/auth"
)

func TestAppendTenantFilter(t *testing.T) {
	query, args, err := AppendTenantFilter("SELECT * FROM users", "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if query != "SELECT * FROM users WHERE tenant_id = ?" {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(args) != 1 || args[0] != "t1" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestAppendTenantFilterMissingTenant(t *testing.T) {
	_, _, err := AppendTenantFilter("SELECT * FROM users", "")
	if !errors.Is(err, auth.ErrMissingTenant) {
		t.Fatalf("expected missing tenant error, got: %v", err)
	}
}
