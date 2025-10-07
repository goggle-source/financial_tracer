package postgresql

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"gorm.io/gorm"
)

func (d *Db) CreateCategoryDatabase(userID uint, category domain.CategoryInput) (uint, error) {
	var user User
	result := d.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, ErrorNotFound
		}
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, ErrorDuplicated
		}
		return 0, result.Error
	}

	newCategory := Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Description: category.Description,
	}

	d.DB.Model(&user).Association("Categories").Append(&newCategory)
	return newCategory.ID, nil
}

func (d *Db) ReadCategoryDatabase(idCategory uint) (domain.CategoryOutput, error) {

	var category Category
	result := d.DB.First(&category, idCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.CategoryOutput{}, ErrorNotFound
		}
		return domain.CategoryOutput{}, result.Error
	}

	modelCategory := domain.CategoryOutput{
		Name:        category.Name,
		Description: category.Description,
		Limit:       category.Limit,
	}
	return modelCategory, nil
}

func (d *Db) UpdateCategoryDatabase(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {

	var categor Category
	result := d.DB.First(&categor, idCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.CategoryOutput{}, ErrorNotFound
		}
		return domain.CategoryOutput{}, result.Error
	}
	categor.Name = newCategory.Name
	categor.Limit = newCategory.Limit
	categor.Description = newCategory.Description

	result = d.DB.Save(&categor)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.CategoryOutput{}, ErrorDuplicated
		}
		return domain.CategoryOutput{}, result.Error
	}

	ResponseCategory := domain.CategoryOutput{
		Name:        categor.Name,
		Description: categor.Description,
		Limit:       categor.Limit,
	}
	return ResponseCategory, ErrorNotFound
}

func (d *Db) DeleteCategoryDatabase(idCategory uint) error {

	result := d.DB.Delete(&Category{}, idCategory)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}
