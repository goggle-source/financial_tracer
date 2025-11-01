package category

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CreateCategoryRepository interface {
	CreateCategory(id uint, category domain.CategoryInput) (uint, error)
}

type GetCategoryRepository interface {
	GetCategory(idCategory uint) (domain.CategoryOutput, error)
}

type UpdateCategoryRepository interface {
	UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
}

type DeleteCategoryRepository interface {
	DeleteCategory(idCategory uint) error
}

type CategoryTypeRepository interface {
	CategoriesType(typeFound string) ([]domain.CategoryOutput, error)
}

type CategoryServer struct {
	d   DeleteCategoryRepository
	c   CreateCategoryRepository
	g   GetCategoryRepository
	u   UpdateCategoryRepository
	t   CategoryTypeRepository
	log *logrus.Logger
}

func CreateCategoryServer(d DeleteCategoryRepository,
	c CreateCategoryRepository,
	u UpdateCategoryRepository,
	g GetCategoryRepository,
	t CategoryTypeRepository,
	log *logrus.Logger) *CategoryServer {
	return &CategoryServer{
		d:   d,
		c:   c,
		g:   g,
		u:   u,
		t:   t,
		log: log,
	}
}

func (cs *CategoryServer) CreateCategory(userID uint, category domain.CategoryInput) (uint, error) {
	const op = "category.CreateCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":      op,
		"user_id": userID,
	})

	log.Info("start create category")

	if err := validator.New().Struct(category); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return 0, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	id, err := cs.c.CreateCategory(userID, category)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("category is not found")

			return 0, fmt.Errorf("%s category is not found: %w", op, ErrNoFound)
		}
		if errors.Is(err, postgresql.ErrorDuplicated) {
			log.WithField("err", err).Error("field duiplicated category")

			return 0, fmt.Errorf("%s field duiplicated category: %w", op, ErrDuplicated)
		}
		log.WithField("err", err).Error("invalid create category")

		return 0, fmt.Errorf("%s invalid create category: %w", op, ErrDatabase)
	}

	log.Info("success create category")

	return id, nil
}

func (cs *CategoryServer) GetCategory(idCategory uint) (domain.CategoryOutput, error) {
	const op = "category.GetCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":          op,
		"category_id": idCategory,
	})

	category, err := cs.g.GetCategory(idCategory)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("category is not found")

			return domain.CategoryOutput{}, fmt.Errorf("%s category is not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("invalid get category")

		return domain.CategoryOutput{}, fmt.Errorf("%s invalid get category: %w", op, ErrDatabase)
	}
	log.Info("success get category")

	return category, nil
}

func (cs *CategoryServer) UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error) {
	const op = "category.UpdateCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":          op,
		"category_id": idCategory,
	})

	log.Info("start update category")

	if err := validator.New().Struct(newCategory); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return domain.CategoryOutput{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	category, err := cs.u.UpdateCategory(idCategory, newCategory)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("category is not found")

			return domain.CategoryOutput{}, fmt.Errorf("%s category is not found: %w", op, ErrNoFound)
		}
		if errors.Is(err, postgresql.ErrorDuplicated) {
			log.WithField("err", err).Error("field duiplicated category")

			return domain.CategoryOutput{}, fmt.Errorf("%s field duiplicated category: %w", op, ErrDuplicated)
		}
		log.WithField("err", err).Error("invalid update category")

		return domain.CategoryOutput{}, fmt.Errorf("%s invalid update category: %w", op, ErrDatabase)
	}
	log.Info("success update category")

	return category, nil
}

func (cs *CategoryServer) DeleteCategory(idCategory uint) error {
	const op = "category.DeleteCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":          op,
		"category_id": idCategory,
	})

	log.Info("start delete category")

	err := cs.d.DeleteCategory(idCategory)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("category is not found")

			return fmt.Errorf("%s category is not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("invalid delete user")

		return fmt.Errorf("%s invalid delete user: %w", op, ErrDatabase)
	}

	log.Info("success delete category")

	return nil
}

func (cs *CategoryServer) CategoryType(typeFound string) ([]domain.CategoryOutput, error) {
	const op = "category.CategoryType"

	if typeFound == "" {
		return []domain.CategoryOutput{}, fmt.Errorf("%s: %w", op, ErrValidateType)
	}

	result, err := cs.t.CategoriesType(typeFound)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			return []domain.CategoryOutput{}, fmt.Errorf("%s: %w", op, ErrNoFound)
		}

		return []domain.CategoryOutput{}, fmt.Errorf("%s: %w", op, ErrDatabase)
	}

	return result, nil
}
