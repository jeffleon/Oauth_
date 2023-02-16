package infraestructure

import (
	"context"

	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(ctx context.Context, client *redis.Client) domain.RedisRepository {
	return &redisRepository{
		client,
		ctx,
	}
}

func (r *redisRepository) HSet(hashKey, field, value string) error {
	return r.client.HSet(r.ctx, hashKey, field, value).Err()
}

func (r *redisRepository) HGet(hashKey, field string) (string, error) {
	return r.client.HGet(r.ctx, hashKey, field).Result()
}
