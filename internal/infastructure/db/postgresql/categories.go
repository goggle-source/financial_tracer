package postgresql

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"gorm.io/gorm"
)

func (d *Db) CreateCategory(id uint, category domain.Category) (uint, error) {
	var user User

	result := d.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, domain.ErrorNotFound
		}
		return 0, result.Error
	}

	newCategory := Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Description: category.Description,
	}

	d.DB.Model(&user).Association("Categories").Append(&newCategory)
	return user.ID, nil
}

func (d *Db) ReadCategory(idUser uint, idCategory uint) (domain.Category, error) {

	var user User

	result := d.DB.Preload("categories").First(&user, idUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Category{}, domain.ErrorNotFound
		}
		return domain.Category{}, result.Error
	}

	for _, category := range user.Categories {
		if category.ID == idCategory {
			return domain.Category{
				Name:        category.Name,
				Limit:       category.Limit,
				Description: category.Description,
			}, nil
		}
	}
	return domain.Category{}, domain.ErrorNotFound
}

func (d *Db) UpdateCategory(idUser uint, idCategory uint, newCategory domain.Category) (uint, error) {

	var user User

	result := d.DB.Preload("categories").First(&user, idUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, domain.ErrorNotFound
		}
		return 0, result.Error
	}

	for _, category := range user.Categories {
		if category.ID == idCategory {
			category = Category{
				Name:        newCategory.Name,
				Limit:       newCategory.Limit,
				Description: category.Description,
			}
			return category.ID, nil
		}
	}
	return 0, domain.ErrorNotFound
}

func (d *Db) DeleteCategory(idUser uint, idCategory uint) error {

	result := d.DB.Where("user_id = ?", idUser).Delete(&Category{}, idCategory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrorNotFound
		}
		return result.Error
	}

	return nil
}
