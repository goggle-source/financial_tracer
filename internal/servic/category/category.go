package category

import (
	"context"
	"strconv"

	"github.com/financial_tracer/internal/domain"
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

type Redis interface {
	HsetCategory(ctx context.Context, id uint, category domain.CategoryOutput) error
	HgetCategory(ctx context.Context, id uint) (map[string]string, error)
	HdelCategory(ctx context.Context, id uint) error
}

type CategoryServer struct {
	d        DeleteCategoryRepository
	c        CreateCategoryRepository
	g        GetCategoryRepository
	u        UpdateCategoryRepository
	t        CategoryTypeRepository
	log      *logrus.Logger
	rbd      Redis
	validate validator.Validate
}

func CreateCategoryServer(d DeleteCategoryRepository,
	c CreateCategoryRepository,
	u UpdateCategoryRepository,
	g GetCategoryRepository,
	t CategoryTypeRepository,
	log *logrus.Logger,
	rbd Redis) *CategoryServer {
	return &CategoryServer{
		d:        d,
		c:        c,
		g:        g,
		u:        u,
		t:        t,
		log:      log,
		rbd:      rbd,
		validate: *validator.New(),
	}
}

func (cs *CategoryServer) CreateCategory(userID uint, category domain.CategoryInput) (uint, error) {
	const op = "category.CreateCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":      op,
		"user_id": userID,
	})

	log.Info("start create category")

	if err := cs.validate.Struct(category); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return 0, err
	}

	id, err := cs.c.CreateCategory(userID, category)
	if err != nil {
		log.Error("error create category: ", err)
		return 0, RegsiterErrorDatabase(err)
	}

	canal := make(chan error)

	go func(canal chan error) {
		category := domain.CategoryOutput{
			Name:        category.Name,
			Description: category.Description,
			Type:        category.Type,
			Limit:       category.Limit,
			UserID:      userID,
		}
		err := cs.rbd.HsetCategory(context.Background(), id, category)
		canal <- err
	}(canal)

	log.Info("success create category")

	if err := <-canal; err != nil {
		log.Error("invalid set in redis", err)
	}

	return id, nil
}

func (cs *CategoryServer) GetCategory(idCategory uint) (domain.CategoryOutput, error) {
	const op = "category.GetCategory"

	log := cs.log.WithFields(logrus.Fields{
		"op":          op,
		"category_id": idCategory,
	})

	log.Info("start get category")

	result, err := cs.rbd.HgetCategory(context.Background(), idCategory)
	if err == nil {
		log.Info("get cateogry is cash: ", result)
		limit, _ := strconv.Atoi(result["limit"])
		userID, _ := strconv.Atoi(result["userID"])
		return domain.CategoryOutput{
			UserID:      uint(userID),
			Name:        result["name"],
			Limit:       limit,
			Description: result["description"],
			Type:        result["type"],
		}, nil
	} else {
		log.Info("error cash", err)
	}

	category, err := cs.g.GetCategory(idCategory)
	if err != nil {
		log.Error("error get category: ", err)
		return domain.CategoryOutput{}, RegsiterErrorDatabase(err)
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

	if err := cs.validate.Struct(newCategory); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return domain.CategoryOutput{}, err
	}

	category, err := cs.u.UpdateCategory(idCategory, newCategory)
	if err != nil {
		log.Error("error update category: ", err)
		return domain.CategoryOutput{}, RegsiterErrorDatabase(err)
	}

	canal := make(chan error, 1)

	go func(chan error) {
		err := cs.rbd.HsetCategory(context.Background(), idCategory, category)
		canal <- err
	}(canal)

	if err := <-canal; err != nil {
		log.Error("invalid set in cash", err)
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
		log.Error("error delete category: ", err)
		return RegsiterErrorDatabase(err)
	}

	canal := make(chan error)

	go func(canal chan error) {
		err := cs.rbd.HdelCategory(context.Background(), idCategory)
		canal <- err
	}(canal)

	if err := <-canal; err != nil {
		log.Error("invalid set in cash", err)
	}

	log.Info("success delete category")

	return nil
}

func (cs *CategoryServer) CategoryType(typeFound string) ([]domain.CategoryOutput, error) {
	const op = "category.CategoryType"

	log := cs.log.WithFields(logrus.Fields{
		"op":   op,
		"type": typeFound,
	})

	log.Info("start gets categories type")

	if typeFound == "" {
		log.Error("invalid type")
		return []domain.CategoryOutput{}, ErrValidateType
	}

	result, err := cs.t.CategoriesType(typeFound)
	if err != nil {
		log.Error("error find category type: ", err)
		return []domain.CategoryOutput{}, RegsiterErrorDatabase(err)
	}

	log.Info("gets categories type")

	return result, nil
}
