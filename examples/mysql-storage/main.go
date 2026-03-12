// MySQL storage example.
// Run: MYSQL_DSN="user:pass@tcp(localhost:3306)/parevo?parseTime=true" go run .
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/permission"
	mysqlstorage "github.com/parevo/core/storage/mysql"
	"github.com/parevo/core/tenant"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root@tcp(localhost:3306)/parevo?parseTime=true"
		log.Printf("Using default MYSQL_DSN. Set MYSQL_DSN for custom connection.")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("MySQL open: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		log.Fatalf("MySQL ping: %v (ensure MySQL is running and schema is applied: mysql < storage/mysql/schema.sql)", err)
	}

	tenantStore := mysqlstorage.NewTenantStore(db)
	permissionStore := mysqlstorage.NewPermissionStore(db)

	// Seed sample data (optional - run once)
	_, _ = db.Exec(`INSERT IGNORE INTO parevo_subject_tenants (subject_id, tenant_id) VALUES ('user-1', 'tenant-a'), ('user-1', 'tenant-b')`)
	_, _ = db.Exec(`INSERT IGNORE INTO parevo_permission_grants (subject_id, tenant_id, permission) VALUES ('user-1', 'tenant-a', 'orders:read')`)

	svc, err := auth.NewServiceWithModules(auth.Config{
		Issuer:    "parevo",
		Audience:  "parevo-api",
		SecretKey: []byte("change-this-in-production"),
	}, auth.Modules{
		Tenant:         tenant.NewService(tenantStore),
		Permission:     permission.NewService(permissionStore),
		TenantOverride: tenant.StaticOverridePolicy{Allow: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	protected := nethttpadapter.AuthMiddleware(svc, adapters.Options{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := auth.ClaimsFromContext(r.Context())
		_, _ = fmt.Fprintf(w, "hello user=%s tenant=%s", claims.UserID, claims.TenantID)
	}))
	mux.Handle("/secure", protected)

	fmt.Println("MySQL storage example: http://localhost:8083/secure")
	fmt.Println("Get a JWT and: curl -H 'Authorization: Bearer <token>' http://localhost:8083/secure")
	_ = http.ListenAndServe(":8083", mux)
}
