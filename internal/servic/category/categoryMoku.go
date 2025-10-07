package category

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) CreateCategoryDatabase(id uint, category domain.CategoryInput) (uint, error) {
	args := d.Called(id, category)
	return args.Get(0).(uint), args.Error(1)
}

func (d *DbMock) ReadCategoryDatabase(idCategory uint) (domain.CategoryOutput, error) {
	args := d.Called(idCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) UpdateCategoryDatabase(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	args := d.Called(idCategory, newCategory)
	return args.Get(0).(domain.CategoryOutput), args.Error(1)
}

func (d *DbMock) DeleteCategoryDatabase(idCategory uint) error {
	args := d.Called(idCategory)
	return args.Error(0)
}
