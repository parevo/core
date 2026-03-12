package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/parevo/core/storage"
)

const (
	sessionKeyPrefix   = "parevo:session:"
	sessionRevokedKey  = "parevo:session:revoked:"
	userSessionsPrefix = "parevo:user:sessions:"
)

type RedisSessionStore struct {
	client *redis.Client
	prefix string
}

func NewSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client, prefix: "parevo"}
}

func (s *RedisSessionStore) RevokeSession(ctx context.Context, sessionID string) error {
	return s.client.Set(ctx, sessionRevokedKey+sessionID, "1", 0).Err()
}

func (s *RedisSessionStore) IsSessionRevoked(ctx context.Context, sessionID string) (bool, error) {
	val, err := s.client.Get(ctx, sessionRevokedKey+sessionID).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "1", err
}

func (s *RedisSessionStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	key := userSessionsPrefix + userID
	if err := s.client.SAdd(ctx, key, sessionID).Err(); err != nil {
		return err
	}
	return s.client.Set(ctx, sessionKeyPrefix+sessionID, userID, 0).Err()
}

func (s *RedisSessionStore) RevokeAllSessionsByUser(ctx context.Context, userID string) error {
	members, err := s.client.SMembers(ctx, userSessionsPrefix+userID).Result()
	if err != nil {
		return err
	}
	for _, sid := range members {
		if err := s.client.Set(ctx, sessionRevokedKey+sid, "1", 0).Err(); err != nil {
			return err
		}
	}
	return s.client.Del(ctx, userSessionsPrefix+userID).Err()
}

var _ storage.UserSessionStore = (*RedisSessionStore)(nil)
