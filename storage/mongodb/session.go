package mongodb

import (
	"context"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SessionStore implements storage.UserSessionStore for MongoDB.
type SessionStore struct {
	coll *mongo.Collection
}

// NewSessionStore creates a MongoDB SessionStore.
func NewSessionStore(db *mongo.Database) *SessionStore {
	return &SessionStore{coll: db.Collection("parevo_sessions")}
}

// RevokeSession marks the session as revoked.
func (s *SessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"session_id": sessionID}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}

// IsSessionRevoked checks if the session is revoked.
func (s *SessionStore) IsSessionRevoked(ctx context.Context, sessionID string) (bool, error) {
	var doc struct {
		Revoked bool `bson:"revoked"`
	}
	err := s.coll.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return doc.Revoked, nil
}

// BindSessionToUser binds a session to a user.
func (s *SessionStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	_, err := s.coll.UpdateOne(ctx,
		bson.M{"session_id": sessionID},
		bson.M{"$set": bson.M{"session_id": sessionID, "user_id": userID, "revoked": false}},
		options.UpdateOne().SetUpsert(true))
	return err
}

// RevokeAllSessionsByUser revokes all sessions for the user.
func (s *SessionStore) RevokeAllSessionsByUser(ctx context.Context, userID string) error {
	_, err := s.coll.UpdateMany(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"revoked": true}})
	return err
}

// ListSessionsByUser returns session IDs for the user (for admin UI).
func (s *SessionStore) ListSessionsByUser(ctx context.Context, userID string) ([]string, error) {
	cur, err := s.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer func() { _ = cur.Close(ctx) }()
	var out []string
	for cur.Next(ctx) {
		var doc struct {
			SessionID string `bson:"session_id"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, doc.SessionID)
	}
	return out, cur.Err()
}

var _ storage.UserSessionStore = (*SessionStore)(nil)
var _ storage.SessionListStore = (*SessionStore)(nil)
