// MongoDB storage example.
// Run: MONGODB_URI="mongodb://localhost:27017" go run .
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/parevo/core/auth"
	"github.com/parevo/core/auth/adapters"
	nethttpadapter "github.com/parevo/core/auth/adapters/nethttp"
	"github.com/parevo/core/permission"
	mongostorage "github.com/parevo/core/storage/mongodb"
	"github.com/parevo/core/tenant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
		log.Printf("Using default MONGODB_URI. Set MONGODB_URI for custom connection.")
	}

	ctx := context.Background()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connect: %v", err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping: %v (ensure MongoDB is running)", err)
	}

	db := client.Database("parevo_example")

	tenantStore := mongostorage.NewTenantStore(db)
	permissionStore := mongostorage.NewPermissionStore(db)

	// Seed sample data (collections auto-created)
	permStoreImpl := mongostorage.NewPermissionStore(db)
	_ = permStoreImpl.Grant(ctx, "user-1", "tenant-a", "orders:read")

	coll := db.Collection("parevo_subject_tenants")
	_, _ = coll.UpdateOne(ctx, bson.M{"subject_id": "user-1", "tenant_id": "tenant-a"},
		bson.M{"$set": bson.M{"subject_id": "user-1", "tenant_id": "tenant-a"}}, options.UpdateOne().SetUpsert(true))
	_, _ = coll.UpdateOne(ctx, bson.M{"subject_id": "user-1", "tenant_id": "tenant-b"},
		bson.M{"$set": bson.M{"subject_id": "user-1", "tenant_id": "tenant-b"}}, options.UpdateOne().SetUpsert(true))

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

	fmt.Println("MongoDB storage example: http://localhost:8084/secure")
	fmt.Println("Get a JWT and: curl -H 'Authorization: Bearer <token>' http://localhost:8084/secure")
	_ = http.ListenAndServe(":8084", mux)
}
