package cash

import (
	"context"

	"github.com/financial_tracer/internal/config"
	"github.com/financial_tracer/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RealRedis struct {
	r *redis.Client
}

func CreateRealRedis(cfg config.Config) RealRedis {
	rbd := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return RealRedis{
		r: rbd,
	}
}

func (r *RealRedis) HsetCategory(ctx context.Context, key string, category domain.CategoryOutput) error {
	return r.r.HSet(context.Background(), key,
		"userID", category.UserID,
		"name", category.Name,
		"type", category.Type,
		"description", category.Description,
		"limit", category.Limit).Err()
}

func (r *RealRedis) HgetCategory(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.r.HGetAll(ctx, key).Result()
	if err != nil {
		return map[string]string{}, err
	}

	return result, nil
}

func (r *RealRedis) HdelCategory(ctx context.Context, key string) error {
	return r.r.HDel(ctx, key).Err()
}
