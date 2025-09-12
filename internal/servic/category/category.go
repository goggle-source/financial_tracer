package category

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/go-playground/validator/v10"
)

type DatabaseCategoryRepository interface {
	CreateCategory(id uint, category domain.Category) (uint, error)
	ReadCategory(idUser uint, idCategory uint) (domain.Category, error)
	UpdateCategory(idUser uint, idCategory uint, newCategory domain.Category) (uint, error)
	DeleteCategory(idUser uint, idCategory uint) error
}

type CategoryServer struct {
	d DatabaseCategoryRepository
}

func CreateCategoryServer(d DatabaseCategoryRepository) *CategoryServer {
	return &CategoryServer{
		d: d,
	}
}

func (cs *CategoryServer) Create(idUser uint, category models.Category) (uint, error) {

	if err := validator.New().Struct(category); err != nil {
		return 0, err
	}

	newCategory := ModelsToDomain(category)

	id, err := cs.d.CreateCategory(idUser, newCategory)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cs *CategoryServer) Get(idUser uint, idCategory uint) (models.Category, error) {

	category, err := cs.d.ReadCategory(idUser, idCategory)
	if err != nil {
		return models.Category{}, err
	}

	ResultCategory := DomainToModels(category)
	return ResultCategory, nil
}

func (cs *CategoryServer) Update(idUser uint, idCategory uint, newCategory models.Category) (uint, error) {

	if err := validator.New().Struct(newCategory); err != nil {
		return 0, err
	}

	cat := ModelsToDomain(newCategory)

	category, err := cs.d.UpdateCategory(idUser, idCategory, cat)
	if err != nil {
		return 0, err
	}

	return category, nil
}

func (cs *CategoryServer) Delete(idUser uint, idCategory uint) error {

	err := cs.d.DeleteCategory(idUser, idCategory)
	if err != nil {
		return err
	}

	return nil
}
