package categoryHandlers

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type categoryServiceMock struct {
	mock.Mock
}

func (m *categoryServiceMock) CreateCategory(idUser uint, category domain.CategoryInput) (uint, error) {
	args := m.Called(idUser, category)
	return args.Get(0).(uint), args.Error(1)
}
func (m *categoryServiceMock) GetCategory(idCategory uint) (domain.CategoryOutput, error) {
	args := m.Called(idCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}
func (m *categoryServiceMock) UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	args := m.Called(idCategory, newCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}
func (m *categoryServiceMock) DeleteCategory(idCategory uint) error {
	args := m.Called(idCategory)
	return args.Error(0)
}
