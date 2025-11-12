package categoryHandlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/handlers/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServicCategoryer interface {
	CreateCategory(ctx context.Context, idUser uint, category domain.CategoryInput) (uint, error)
	ReadCategory(ctx context.Context, idCategory uint) (domain.CategoryOutput, error)
	UpdateCategory(ctx context.Context, idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
	DeleteCategory(ctx context.Context, idCategory uint) error
}

type CreateCategoryServic interface {
	CreateCategory(ctx context.Context, idUser uint, category domain.CategoryInput) (uint, error)
}

type GetCategoryServic interface {
	GetCategory(ctx context.Context, idCategory uint) (domain.CategoryOutput, error)
}

type UpdateCategoryServic interface {
	UpdateCategory(ctx context.Context, idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
}

type DeleteCategoryServic interface {
	DeleteCategory(ctx context.Context, idCategory uint) error
}

type CategoryTypeServic interface {
	CategoryType(ctx context.Context, typeFound string) ([]domain.CategoryOutput, error)
}

type CategoryHandlers struct {
	c   CreateCategoryServic
	g   GetCategoryServic
	u   UpdateCategoryServic
	d   DeleteCategoryServic
	t   CategoryTypeServic
	log *logrus.Logger
	ctx context.Context
}

func CreateHandlersCategory(c CreateCategoryServic,
	g GetCategoryServic,
	u UpdateCategoryServic,
	d DeleteCategoryServic,
	t CategoryTypeServic,
	log *logrus.Logger,
	ctx context.Context) *CategoryHandlers {
	return &CategoryHandlers{
		c:   c,
		g:   g,
		u:   u,
		d:   d,
		t:   t,
		log: log,
		ctx: ctx,
	}
}

// CreateCategory godoc
//
//	@Summary		Создание Категории
//	@Description	Создание новой категории для зарегистрировшегося пользователя
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//
//	@Param			req	body		RequestCreateCategory	true	"данные для создание новой категории"
//
//	@Success		200	{object}	api.SuccessResponse		"Успешное создание категории"
//
//	@Failure		401	{object}	api.ErrorResponse		"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse		"Некорректные входные данные"
//	@Failure		400	{object}	api.ErrorResponse		"Некорректный данные"
//
//	@Failure		501	{object}	api.ErrorResponse		"Ошибка сервера"
//
//	@Router			/category/ [post]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) PostCategory(c *gin.Context) {
	const op = "handlers.PostCategory"

	log := h.log.WithField("op", op)

	log.Info("start create category")

	var newCategory RequestCreateCategory
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		log.WithField("err", err).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	idUser, ok := c.Get("userID")
	h.log.Info(idUser)
	if !ok {
		log.Error("error get userID")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	cat := domain.CategoryInput{
		Name:        newCategory.Name,
		Limit:       newCategory.Limit,
		Type:        newCategory.Type,
		Description: newCategory.Description,
	}

	idCategory, err := h.c.CreateCategory(h.ctx, idUser.(uint), cat)
	if err != nil {
		log.Error("error create category")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success create category")

	api.ResponseOK(c, idCategory)
}

// GetCategory godoc
//
//	@Summary		Получение категории
//	@Description	Получение категории для конкретного пользователя
//	@Tags			categories
//	@Produce		json
//	@Param			id	path		int					true	"ID категории"
//	@Success		200	{object}	api.SuccessResponse	"success"
//
//	@Failure		401	{object}	api.ErrorResponse	"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse	"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse	"Некорректный данные"
//	@Failure		404	{object}	api.ErrorResponse	"Категория не найдена"
//
//	@Router			/category/{id} [get]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) GetCategory(c *gin.Context) {
	const op = "handlers.GetCategory"

	log := h.log.WithField("op", op)

	log.Info("start get cateogry")

	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.WithField("err", err).Error("error get id category")
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.g.GetCategory(h.ctx, uint(id))
	if err != nil {
		log.Error("error get category")
		api.RegistrationError(c, err)
		return
	}
	api.ResponseOK(c, category)
}

// UpdateCategory godoc
//
//	@Summary		обновление пользователя
//	@Description	Обновление всех характеристики категории
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestUpdateCategory	true	"данные для обновление категории"
//	@Success		200	{object}	api.SuccessResponse		"success"
//
//	@Failure		401	{object}	api.ErrorResponse		"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse		"Некоректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse		"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse		"Некорректный данные"
//	@Failure		404	{object}	api.ErrorResponse		"Категория не найдена"
//
//	@Router			/category/ [put]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) UpdateCategory(c *gin.Context) {
	const op = "handlers.UpdateCategory"

	log := h.log.WithField("op", op)

	log.Info("start update category")

	var updateCategory RequestUpdateCategory

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		log.WithField("err", err).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	newCategory := domain.CategoryInput{
		Name:        updateCategory.Name,
		Limit:       updateCategory.Limit,
		Type:        updateCategory.Type,
		Description: updateCategory.Description,
	}

	category, err := h.u.UpdateCategory(h.ctx, updateCategory.CategoryId, newCategory)
	if err != nil {
		log.Error("error update category")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success update category")

	api.ResponseOK(c, category)
}

// DeleteCategory godoc
//
//	@Summary		Удаление категории
//	@Description	Удаление категории конкретного пользователя
//	@Tags			categories
//	@Produce		json
//	@Param			id	path		int					true	"ID категории"
//	@Success		200	{object}	api.SuccessResponse	"success"
//
//	@Failure		401	{object}	api.ErrorResponse	"Ошибка авторизации"
//
//	@Failure		500	{object}	api.ErrorResponse	"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse	"Некоректные входные данные"
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse	"Категория не найдена"
//
//	@Router			/category/{id} [delete]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) DeleteCategory(c *gin.Context) {
	const op = "handlers.DeleteCategory"

	log := h.log.WithField("op", op)

	log.Info("start delete category")

	param := c.Param("id")
	id, err := strconv.Atoi(param)

	if err != nil {
		log.WithField("err", err).Error("error get id")
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.d.DeleteCategory(h.ctx, uint(id))
	if err != nil {
		log.Error("error delete category")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success delete category")

	api.ResponseOK(c, "user delete")
}

// CategoryType godoc
// @Summary получение категории или категорий по типу
// @Description получение категории или категорий по определенномк типу, который передается в URL пути
// @Tags categories
// @Produce json
// @Param type path string true "тип, по которому будут выбиратся категории или категория"
// @Success 200 {object}  api.SuccessResponse "success"
//
//	@Failure	500	{object}	api.ErrorResponse	"Ошибка сервера"
//	@Failure	400	{object}	api.ErrorResponse	"Некоректные входные данные"
//	@Failure	400	{object}	api.ErrorResponse	"Некорректные данные"
//	@Failure	404	{object}	api.ErrorResponse	"Категория не найдена"
//	@Router		/category/type/{type} [get]
//	@Security	jwtAuth
func (h *CategoryHandlers) CategoryType(c *gin.Context) {
	const op = "handlers.CategoryType"

	log := h.log.WithField("op", op)

	log.Info("start gets categories type")

	typeFound := c.Param("type")
	if typeFound == "" {
		log.WithField("err", "invalid param").Error("param is not valid")
		api.ResponseError(c, http.StatusBadRequest, "invalid param")
		return
	}

	result, err := h.t.CategoryType(h.ctx, typeFound)
	if err != nil {
		log.WithField("err", err).Error("error in get category type")
		api.RegistrationError(c, err)
	}

	log.Info("success get category type")

	api.ResponseOK(c, result)
}
