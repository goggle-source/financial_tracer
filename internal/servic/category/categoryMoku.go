package category

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) CreateCategory(id uint, category domain.CategoryInput) (uint, error) {
	args := d.Called(id, category)
	return args.Get(0).(uint), args.Error(1)
}

func (d *DbMock) GetCategory(idCategory uint) (domain.CategoryOutput, error) {
	args := d.Called(idCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	args := d.Called(idCategory, newCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) DeleteCategory(idCategory uint) error {
	args := d.Called(idCategory)
	return args.Error(0)
}

func (d *DbMock) CategoriesType(typeFound string) ([]domain.CategoryOutput, error) {
	args := d.Called(typeFound)
	return args.Get(0).([]domain.CategoryOutput), args.Error(1)
}
