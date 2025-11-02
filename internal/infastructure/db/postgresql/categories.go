package postgresql

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"gorm.io/gorm"
)

func (d *Db) CreateCategory(userID uint, category domain.CategoryInput) (uint, error) {
	var user User
	result := d.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, ErrorNotFound
		}
		return 0, result.Error
	}

	newCategory := Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Description: category.Description,
	}

	err := d.DB.Model(&user).Association("Categories").Append(&newCategory)
	if err != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, ErrorDuplicated
		}
		return 0, err
	}
	return newCategory.ID, nil
}

func (d *Db) GetCategory(idCategory uint) (domain.CategoryOutput, error) {

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

func (d *Db) UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {

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

func (d *Db) DeleteCategory(idCategory uint) error {

	result := d.DB.Delete(&Category{}, idCategory)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}

func (d *Db) CategoriesType(typeFound string) ([]domain.CategoryOutput, error) {

	var category []Category

	result := d.DB.Find(&category, "type = ?", category)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []domain.CategoryOutput{}, ErrorNotFound
		}

		return []domain.CategoryOutput{}, result.Error
	}

	var arr []domain.CategoryOutput

	for _, value := range category {
		categor := domain.CategoryOutput{
			UserID:      value.UserID,
			Name:        value.Name,
			Limit:       value.Limit,
			Type:        value.Type,
			Description: value.Description,
		}

		arr = append(arr, categor)
	}

	return arr, nil

}
