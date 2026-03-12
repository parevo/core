package mongodb

import (
	"context"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TenantStore implements storage.TenantStore for MongoDB.
type TenantStore struct {
	coll *mongo.Collection
}

// NewTenantStore creates a MongoDB TenantStore.
func NewTenantStore(db *mongo.Database) *TenantStore {
	return &TenantStore{coll: db.Collection("parevo_subject_tenants")}
}

// ResolveSubjectTenants returns tenant IDs the subject has access to.
func (s *TenantStore) ResolveSubjectTenants(ctx context.Context, subjectID string) ([]string, error) {
	cur, err := s.coll.Find(ctx, bson.M{"subject_id": subjectID})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var tenants []string
	for cur.Next(ctx) {
		var doc struct {
			TenantID string `bson:"tenant_id"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		tenants = append(tenants, doc.TenantID)
	}
	return tenants, cur.Err()
}

var _ storage.TenantStore = (*TenantStore)(nil)

// TenantLifecycleStore implements storage.TenantLifecycleStore and storage.TenantListStore for MongoDB.
type TenantLifecycleStore struct {
	coll *mongo.Collection
}

// NewTenantLifecycleStore creates a MongoDB TenantLifecycleStore.
func NewTenantLifecycleStore(db *mongo.Database) *TenantLifecycleStore {
	return &TenantLifecycleStore{coll: db.Collection("parevo_tenants")}
}

// Create creates a new tenant.
func (s *TenantLifecycleStore) Create(ctx context.Context, tenantID, name, ownerID string) error {
	_, err := s.coll.UpdateOne(ctx,
		bson.M{"tenant_id": tenantID},
		bson.M{"$setOnInsert": bson.M{"tenant_id": tenantID, "name": name, "owner_id": ownerID, "status": "active"}},
		options.UpdateOne().SetUpsert(true))
	return err
}

// Suspend marks the tenant as suspended.
func (s *TenantLifecycleStore) Suspend(ctx context.Context, tenantID string) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"status": "suspended"}})
	return err
}

// Resume marks the tenant as active.
func (s *TenantLifecycleStore) Resume(ctx context.Context, tenantID string) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"status": "active"}})
	return err
}

// Delete marks the tenant as deleted (soft delete).
func (s *TenantLifecycleStore) Delete(ctx context.Context, tenantID string) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"status": "deleted"}})
	return err
}

// Status returns the tenant status.
func (s *TenantLifecycleStore) Status(ctx context.Context, tenantID string) (storage.TenantStatus, error) {
	var doc struct {
		Status string `bson:"status"`
	}
	err := s.coll.FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return storage.TenantStatus(doc.Status), nil
}

// ListTenants returns all non-deleted tenants.
func (s *TenantLifecycleStore) ListTenants(ctx context.Context) ([]storage.TenantInfo, error) {
	cur, err := s.coll.Find(ctx, bson.M{"status": bson.M{"$ne": "deleted"}})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var out []storage.TenantInfo
	for cur.Next(ctx) {
		var doc struct {
			TenantID string `bson:"tenant_id"`
			Status   string `bson:"status"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, storage.TenantInfo{ID: doc.TenantID, Status: storage.TenantStatus(doc.Status)})
	}
	return out, cur.Err()
}

var _ storage.TenantLifecycleStore = (*TenantLifecycleStore)(nil)
var _ storage.TenantListStore = (*TenantLifecycleStore)(nil)
