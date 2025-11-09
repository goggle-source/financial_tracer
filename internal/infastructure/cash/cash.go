package cash

import (
	"context"
	"strconv"

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

func (r *RealRedis) HsetCategory(ctx context.Context, id uint, category domain.CategoryOutput) error {
	strId := strconv.FormatUint(uint64(id), 10)
	return r.r.HSet(context.Background(), "category:"+strId,
		"userID", category.UserID,
		"name", category.Name,
		"type", category.Type,
		"description", category.Description,
		"limit", category.Limit).Err()
}

func (r *RealRedis) HgetCategory(ctx context.Context, id uint) (map[string]string, error) {
	strId := strconv.FormatUint(uint64(id), 10)
	result, err := r.r.HGetAll(ctx, "category:"+strId).Result()
	if err != nil {
		return map[string]string{}, err
	}

	return result, nil
}

func (r *RealRedis) HdelCategory(ctx context.Context, id uint) error {
	strId := strconv.FormatUint(uint64(id), 10)
	return r.r.HDel(ctx, "category:"+strId).Err()
}
