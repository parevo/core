package mongodb

import (
	"context"
	"time"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RefreshStore implements storage.UserRefreshTokenStore for MongoDB.
type RefreshStore struct {
	coll   *mongo.Collection
	sessColl *mongo.Collection
}

// NewRefreshStore creates a MongoDB RefreshStore.
func NewRefreshStore(db *mongo.Database) *RefreshStore {
	return &RefreshStore{
		coll:     db.Collection("parevo_refresh_tokens"),
		sessColl: db.Collection("parevo_sessions"),
	}
}

// MarkIssued records a refresh token as issued.
func (s *RefreshStore) MarkIssued(ctx context.Context, sessionID, tokenID string, expiresAt time.Time) error {
	var sess struct {
		UserID string `bson:"user_id"`
	}
	err := s.sessColl.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&sess)
	if err != nil {
		return err
	}
	_, err = s.coll.InsertOne(ctx, bson.M{
		"token_id":   tokenID,
		"session_id": sessionID,
		"user_id":    sess.UserID,
		"expires_at": expiresAt,
	})
	return err
}

// IsUsed checks if the token was already used.
func (s *RefreshStore) IsUsed(ctx context.Context, tokenID string) (bool, error) {
	var doc struct {
		ReplacedBy *string `bson:"replaced_by"`
	}
	err := s.coll.FindOne(ctx, bson.M{"token_id": tokenID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return doc.ReplacedBy != nil, nil
}

// MarkUsed marks the token as used.
func (s *RefreshStore) MarkUsed(ctx context.Context, tokenID, replacedBy string) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"token_id": tokenID}, bson.M{"$set": bson.M{"replaced_by": replacedBy}})
	return err
}

// RevokeSession revokes all refresh tokens for the session.
func (s *RefreshStore) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.coll.UpdateMany(ctx, bson.M{"session_id": sessionID}, bson.M{"$set": bson.M{"replaced_by": "revoked"}})
	return err
}

// BindSessionToUser is a no-op for MongoDB.
func (s *RefreshStore) BindSessionToUser(_ context.Context, _, _ string) error {
	return nil
}

// RevokeAllByUser revokes all refresh tokens for the user.
func (s *RefreshStore) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := s.coll.UpdateMany(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"replaced_by": "revoked"}})
	return err
}

var _ storage.UserRefreshTokenStore = (*RefreshStore)(nil)
