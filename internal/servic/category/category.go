package category

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type DatabaseCategoryRepository interface {
	CreateCategoryDatabase(id uint, category domain.CategoryInput) (uint, error)
	ReadCategoryDatabase(idCategory uint) (domain.CategoryOutput, error)
	UpdateCategoryDatabase(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
	DeleteCategoryDatabase(idCategory uint) error
}

type CategoryServer struct {
	d   DatabaseCategoryRepository
	log *logrus.Logger
}

func CreateCategoryServer(d DatabaseCategoryRepository, log *logrus.Logger) *CategoryServer {
	return &CategoryServer{
		d:   d,
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

	id, err := cs.d.CreateCategoryDatabase(userID, category)
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

func (cs *CategoryServer) ReadCategory(idCategory uint) (domain.CategoryOutput, error) {
	const op = "category.ReadCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":          op,
		"category_id": idCategory,
	})

	category, err := cs.d.ReadCategoryDatabase(idCategory)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("category is not found")

			return domain.CategoryOutput{}, fmt.Errorf("%s category is not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("invalid read category")

		return domain.CategoryOutput{}, fmt.Errorf("%s invalid read category: %w", op, ErrDatabase)
	}
	log.Info("success read category")

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

	category, err := cs.d.UpdateCategoryDatabase(idCategory, newCategory)
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

	err := cs.d.DeleteCategoryDatabase(idCategory)
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
