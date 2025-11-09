package cash

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (r *RedisMock) HsetCategory(ctx context.Context, key string, category domain.CategoryOutput) error {
	args := r.Called(ctx, key, category)
	return args.Error(0)
}

func (r *RedisMock) HgetCategory(ctx context.Context, key string) (map[string]string, error) {
	args := r.Called(ctx, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (r *RedisMock) HdelCategory(ctx context.Context, key string) error {
	args := r.Called(ctx, key)
	return args.Error(0)
}
