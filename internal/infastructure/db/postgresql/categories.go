package postgresql

import (
	"context"
	"errors"

	"github.com/financial_tracer/internal/domain"
	"gorm.io/gorm"
)

func (d *Db) CreateCategory(ctx context.Context, userID uint, category domain.CategoryInput) (uint, error) {
	var user User
	result := d.DB.WithContext(ctx).First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, ErrorNotFound
		}
		return 0, result.Error
	}

	newCategory := Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Type:        category.Type,
		Description: category.Description,
	}

	err := d.DB.WithContext(ctx).Model(&user).Association("Categories").Append(&newCategory)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, ErrorDuplicated
		}
		return 0, err
	}
	return newCategory.ID, nil
}

func (d *Db) GetCategory(ctx context.Context, idCategory uint) (domain.CategoryOutput, error) {

	var category Category
	result := d.DB.WithContext(ctx).First(&category, idCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.CategoryOutput{}, ErrorNotFound
		}
		return domain.CategoryOutput{}, result.Error
	}

	modelCategory := domain.CategoryOutput{
		Name:        category.Name,
		Description: category.Description,
		Type:        category.Type,
		Limit:       category.Limit,
	}
	return modelCategory, nil
}

func (d *Db) UpdateCategory(ctx context.Context, idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {

	var categor Category
	result := d.DB.WithContext(ctx).First(&categor, idCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.CategoryOutput{}, ErrorNotFound
		}
		return domain.CategoryOutput{}, result.Error
	}

	result = d.DB.WithContext(ctx).Model(&categor).Updates(Category{Name: newCategory.Name,
		Description: newCategory.Description,
		Type:        newCategory.Type},
	)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.CategoryOutput{}, ErrorDuplicated
		}
		return domain.CategoryOutput{}, result.Error
	}

	ResponseCategory := domain.CategoryOutput{
		Name:        categor.Name,
		Description: categor.Description,
		Type:        categor.Type,
		Limit:       categor.Limit,
	}
	return ResponseCategory, nil
}

func (d *Db) DeleteCategory(ctx context.Context, idCategory uint) error {

	result := d.DB.WithContext(ctx).Unscoped().Delete(&Category{}, idCategory)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *Db) CategoriesType(ctx context.Context, typeFound string) ([]domain.CategoryOutput, error) {

	var category []Category

	result := d.DB.WithContext(ctx).Where("type = ?", typeFound).Find(&category, category)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []domain.CategoryOutput{}, ErrorNotFound
		}

		return []domain.CategoryOutput{}, result.Error
	}

	var arr []domain.CategoryOutput

	for _, value := range category {
		categor := domain.CategoryOutput{
			UserID:      0,
			Name:        value.Name,
			Limit:       value.Limit,
			Type:        value.Type,
			Description: value.Description,
		}

		arr = append(arr, categor)
	}

	return arr, nil

}
