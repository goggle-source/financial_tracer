package category

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) CreateCategory(ctx context.Context, id uint, category domain.CategoryInput) (uint, error) {
	args := d.Called(ctx, id, category)
	return args.Get(0).(uint), args.Error(1)
}

func (d *DbMock) GetCategory(ctx context.Context, idCategory uint) (domain.CategoryOutput, error) {
	args := d.Called(ctx, idCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) UpdateCategory(ctx context.Context, idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	args := d.Called(ctx, idCategory, newCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) DeleteCategory(ctx context.Context, idCategory uint) error {
	args := d.Called(ctx, idCategory)
	return args.Error(0)
}

func (d *DbMock) CategoriesType(ctx context.Context, typeFound string) ([]domain.CategoryOutput, error) {
	args := d.Called(ctx, typeFound)
	return args.Get(0).([]domain.CategoryOutput), args.Error(1)
}
