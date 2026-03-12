package mongodb

import (
	"context"
	"strings"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// PermissionStore implements storage.PermissionStore and storage.PermissionGrantStore for MongoDB.
type PermissionStore struct {
	coll *mongo.Collection
}

// NewPermissionStore creates a MongoDB PermissionStore.
func NewPermissionStore(db *mongo.Database) *PermissionStore {
	return &PermissionStore{coll: db.Collection("parevo_permission_grants")}
}

// HasPermission checks if the subject has the permission in the tenant.
func (s *PermissionStore) HasPermission(ctx context.Context, subjectID, tenantID, permission string, _ []string) (bool, error) {
	cur, err := s.coll.Find(ctx, bson.M{"subject_id": subjectID, "tenant_id": tenantID})
	if err != nil {
		return false, err
	}
	defer func() { _ = cur.Close(ctx) }()

	for cur.Next(ctx) {
		var doc struct {
			Permission string `bson:"permission"`
		}
		if err := cur.Decode(&doc); err != nil {
			return false, err
		}
		if matchPermission(doc.Permission, permission) {
			return true, nil
		}
	}
	return false, cur.Err()
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
	_, err := s.coll.UpdateOne(ctx,
		bson.M{"subject_id": subjectID, "tenant_id": tenantID, "permission": permission},
		bson.M{"$set": bson.M{"subject_id": subjectID, "tenant_id": tenantID, "permission": permission}},
		options.UpdateOne().SetUpsert(true))
	return err
}

// Revoke removes a permission grant.
func (s *PermissionStore) Revoke(ctx context.Context, subjectID, tenantID, permission string) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"subject_id": subjectID, "tenant_id": tenantID, "permission": permission})
	return err
}

// ListGrants returns all permissions for the subject in the tenant.
func (s *PermissionStore) ListGrants(ctx context.Context, subjectID, tenantID string) ([]string, error) {
	cur, err := s.coll.Find(ctx, bson.M{"subject_id": subjectID, "tenant_id": tenantID})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()

	var out []string
	for cur.Next(ctx) {
		var doc struct {
			Permission string `bson:"permission"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, doc.Permission)
	}
	return out, cur.Err()
}

var _ storage.PermissionStore = (*PermissionStore)(nil)
var _ storage.PermissionGrantStore = (*PermissionStore)(nil)
