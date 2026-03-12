package mongodb

import (
	"context"
	"time"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// UserStore implements storage.UserListStore for MongoDB.
type UserStore struct {
	coll *mongo.Collection
}

// NewUserStore creates a MongoDB UserStore.
func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{coll: db.Collection("parevo_users")}
}

// ListUsers returns users from parevo_users (for admin UI).
func (s *UserStore) ListUsers(ctx context.Context) ([]storage.UserInfo, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()
	var out []storage.UserInfo
	for cur.Next(ctx) {
		var doc struct {
			UserID      string    `bson:"user_id"`
			Email       string    `bson:"email"`
			DisplayName string    `bson:"display_name"`
			CreatedAt   time.Time `bson:"created_at"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, storage.UserInfo{
			ID:        doc.UserID,
			Email:     doc.Email,
			Name:      doc.DisplayName,
			CreatedAt: doc.CreatedAt,
		})
	}
	return out, cur.Err()
}

var _ storage.UserListStore = (*UserStore)(nil)
