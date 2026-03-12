package mongodb

import (
	"context"
	"fmt"

	"github.com/parevo/core/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SocialAccountStore implements storage.SocialAccountStore for MongoDB.
type SocialAccountStore struct {
	usersColl *mongo.Collection
	socialColl *mongo.Collection
}

// NewSocialAccountStore creates a MongoDB SocialAccountStore.
func NewSocialAccountStore(db *mongo.Database) *SocialAccountStore {
	return &SocialAccountStore{
		usersColl:  db.Collection("parevo_users"),
		socialColl: db.Collection("parevo_social_accounts"),
	}
}

// FindUserBySocial returns the userID for the given provider and provider user ID.
func (s *SocialAccountStore) FindUserBySocial(ctx context.Context, provider, providerUserID string) (userID string, found bool, err error) {
	var doc struct {
		UserID string `bson:"user_id"`
	}
	err = s.socialColl.FindOne(ctx, bson.M{"provider": provider, "provider_user_id": providerUserID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return doc.UserID, true, nil
}

// FindOrCreateUserByEmail finds or creates a user by email.
func (s *SocialAccountStore) FindOrCreateUserByEmail(ctx context.Context, email, displayName string) (string, error) {
	var doc struct {
		UserID string `bson:"user_id"`
	}
	err := s.usersColl.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err == nil {
		return doc.UserID, nil
	}
	if err != mongo.ErrNoDocuments {
		return "", err
	}
	userID := "user:" + email
	_, err = s.usersColl.InsertOne(ctx, bson.M{
		"user_id":     userID,
		"email":       email,
		"display_name": displayName,
	})
	if err != nil {
		return "", fmt.Errorf("social: create user: %w", err)
	}
	return userID, nil
}

// LinkSocialAccount links a social identity to a user.
func (s *SocialAccountStore) LinkSocialAccount(ctx context.Context, userID string, identity storage.SocialIdentity) error {
	_, err := s.socialColl.UpdateOne(ctx,
		bson.M{"provider": identity.Provider, "provider_user_id": identity.ProviderUserID},
		bson.M{"$set": bson.M{
			"user_id":        userID,
			"email":          identity.Email,
			"email_verified": identity.EmailVerified,
			"name":           identity.Name,
			"avatar_url":     identity.AvatarURL,
		}},
		options.UpdateOne().SetUpsert(true))
	return err
}

var _ storage.SocialAccountStore = (*SocialAccountStore)(nil)
