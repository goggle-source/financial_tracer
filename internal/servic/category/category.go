package category

import (
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/go-playground/validator/v10"
)

type DatabaseCategoryRepository interface {
	CreateCategoryDatabase(id uint, category domain.CategoryInput) (uint, error)
	ReadCategoryDatabase(idCategory uint) (domain.CategoryOutput, error)
	UpdateCategoryDatabase(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
	DeleteCategoryDatabase(idCategory uint) error
}

type CategoryServer struct {
	d DatabaseCategoryRepository
}

func CreateCategoryServer(d DatabaseCategoryRepository) *CategoryServer {
	return &CategoryServer{
		d: d,
	}
}

func (cs *CategoryServer) CreateCategory(userID uint, category domain.CategoryInput) (uint, error) {

	if err := validator.New().Struct(category); err != nil {
		return 0, fmt.Errorf("error validate: %w", err)
	}

	id, err := cs.d.CreateCategoryDatabase(userID, category)
	if err != nil {
		return 0, fmt.Errorf("error create category: %w", err)
	}

	return id, nil
}

func (cs *CategoryServer) ReadCategory(idCategory uint) (domain.CategoryOutput, error) {

	category, err := cs.d.ReadCategoryDatabase(idCategory)
	if err != nil {
		return domain.CategoryOutput{}, fmt.Errorf("error read category: %w", err)
	}

	return category, nil
}

func (cs *CategoryServer) UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {

	if err := validator.New().Struct(newCategory); err != nil {
		return domain.CategoryOutput{}, fmt.Errorf("error validate: %w", err)
	}

	category, err := cs.d.UpdateCategoryDatabase(idCategory, newCategory)
	if err != nil {
		return domain.CategoryOutput{}, fmt.Errorf("error update category: %w", err)
	}

	return category, nil
}

func (cs *CategoryServer) DeleteCategory(idCategory uint) error {

	err := cs.d.DeleteCategoryDatabase(idCategory)
	if err != nil {
		return fmt.Errorf("error delete user: %w", err)
	}

	return nil
}
