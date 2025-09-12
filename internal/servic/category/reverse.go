package category

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
)

func ModelsToDomain(category models.Category) domain.Category {
	return domain.Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Description: category.Description,
	}
}

func DomainToModels(category domain.Category) models.Category {
	return models.Category{
		Name:        category.Name,
		Limit:       category.Limit,
		Description: category.Description,
	}
}
