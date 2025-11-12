package categoryHandlers

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type categoryServiceMock struct {
	mock.Mock
}

func (m *categoryServiceMock) CreateCategory(ctx context.Context, idUser uint, category domain.CategoryInput) (uint, error) {
	args := m.Called(ctx, idUser, category)
	return args.Get(0).(uint), args.Error(1)
}
func (m *categoryServiceMock) GetCategory(ctx context.Context, idCategory uint) (domain.CategoryOutput, error) {
	args := m.Called(ctx, idCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}
func (m *categoryServiceMock) UpdateCategory(ctx context.Context, idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	args := m.Called(ctx, idCategory, newCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}
func (m *categoryServiceMock) DeleteCategory(ctx context.Context, idCategory uint) error {
	args := m.Called(ctx, idCategory)
	return args.Error(0)
}

func (m *categoryServiceMock) CategoryType(ctx context.Context, typeFound string) ([]domain.CategoryOutput, error) {
	args := m.Called(ctx, typeFound)
	return args.Get(0).([]domain.CategoryOutput), args.Error(1)
}
