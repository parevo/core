package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/parevo/core/storage"
)

const (
	refreshUsedPrefix   = "parevo:refresh:used:"
	refreshSessionPrefix = "parevo:refresh:session:"
	refreshUserPrefix   = "parevo:refresh:user:"
)

type RedisRefreshStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRefreshStore(client *redis.Client, defaultTTL time.Duration) *RedisRefreshStore {
	if defaultTTL <= 0 {
		defaultTTL = 7 * 24 * time.Hour
	}
	return &RedisRefreshStore{client: client, ttl: defaultTTL}
}

func (s *RedisRefreshStore) MarkIssued(ctx context.Context, sessionID, tokenID string, expiresAt time.Time) error {
	userID, err := s.client.Get(ctx, sessionKeyPrefix+sessionID).Result()
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	pipe := s.client.Pipeline()
	pipe.Set(ctx, refreshUsedPrefix+tokenID, "issued", time.Until(expiresAt))
	pipe.SAdd(ctx, refreshSessionPrefix+sessionID, tokenID)
	pipe.SAdd(ctx, refreshUserPrefix+userID, sessionID)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisRefreshStore) IsUsed(ctx context.Context, tokenID string) (bool, error) {
	val, err := s.client.Get(ctx, refreshUsedPrefix+tokenID).Result()
	if err == redis.Nil {
		return false, nil
	}
	return val != "" && val != "issued", err
}

func (s *RedisRefreshStore) MarkUsed(ctx context.Context, tokenID, replacedBy string) error {
	return s.client.Set(ctx, refreshUsedPrefix+tokenID, replacedBy, s.ttl).Err()
}

func (s *RedisRefreshStore) RevokeSession(ctx context.Context, sessionID string) error {
	members, err := s.client.SMembers(ctx, refreshSessionPrefix+sessionID).Result()
	if err != nil {
		return err
	}
	pipe := s.client.Pipeline()
	for _, tid := range members {
		pipe.Set(ctx, refreshUsedPrefix+tid, "revoked", s.ttl)
	}
	pipe.Del(ctx, refreshSessionPrefix+sessionID)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisRefreshStore) BindSessionToUser(ctx context.Context, userID, sessionID string) error {
	return s.client.SAdd(ctx, refreshUserPrefix+userID, sessionID).Err()
}

func (s *RedisRefreshStore) RevokeAllByUser(ctx context.Context, userID string) error {
	sessions, err := s.client.SMembers(ctx, refreshUserPrefix+userID).Result()
	if err != nil {
		return err
	}
	for _, sid := range sessions {
		if err := s.RevokeSession(ctx, sid); err != nil {
			return err
		}
	}
	return s.client.Del(ctx, refreshUserPrefix+userID).Err()
}

var _ storage.UserRefreshTokenStore = (*RedisRefreshStore)(nil)
