package categoryHandlers

import (
	"net/http"
	"strconv"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/handlers/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServicCategoryer interface {
	CreateCategory(idUser uint, category domain.CategoryInput) (uint, error)
	ReadCategory(idCategory uint) (domain.CategoryOutput, error)
	UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
	DeleteCategory(idCategory uint) error
}

type CreateCategoryServic interface {
	CreateCategory(idUser uint, category domain.CategoryInput) (uint, error)
}

type GetCategoryServic interface {
	GetCategory(idCategory uint) (domain.CategoryOutput, error)
}

type UpdateCategoryServic interface {
	UpdateCategory(idCategory uint, newCategory domain.CategoryInput) (domain.CategoryOutput, error)
}

type DeleteCategoryServic interface {
	DeleteCategory(idCategory uint) error
}

type CategoryHandlers struct {
	c   CreateCategoryServic
	g   GetCategoryServic
	u   UpdateCategoryServic
	d   DeleteCategoryServic
	log *logrus.Logger
}

func CreateHandlersCategory(c CreateCategoryServic,
	g GetCategoryServic,
	u UpdateCategoryServic,
	d DeleteCategoryServic,
	log *logrus.Logger) *CategoryHandlers {
	return &CategoryHandlers{
		c:   c,
		g:   g,
		u:   u,
		d:   d,
		log: log,
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
//	@Param			req	body		RequestCreateCategory					true	"данные для создание новой категории"
//
//	@Success		200	{object}	api.SuccessResponse[uint]				"Успешное создание категории"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректный данные"
//
//	@Failure		501	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//
//	@Router			/category/create [post]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) PostCategory(c *gin.Context) {
	const op = "handlers.PostCategory"

	c.Writer.Header().Set("content-type", "application/json")
	newCategory := domain.CategoryInput{}
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	idUser, ok := c.Get("userID")
	if !ok {
		h.log.WithField("op", op).Error("error get userID")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	cat := domain.CategoryInput{
		Name:        newCategory.Name,
		Limit:       newCategory.Limit,
		Description: newCategory.Description,
	}

	idCategory, err := h.c.CreateCategory(idUser.(uint), cat)
	if err != nil {
		h.log.WithField("op", op).Error("error create category")
		api.RegistrationError(c, err)
		return
	}
	api.ResponseOK(c, idCategory)
}

// GetCategory godoc
//
//	@Summary		Получение категории
//	@Description	Получение категории для конкретного пользователя
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int								true	"userID для получение категории"
//	@Success		200	{object}	api.SuccessResponse[domain.Category]	"success"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректный данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Категория не найдена"
//
//	@Router			/category/get [get]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) GetCategory(c *gin.Context) {
	const op = "handlers.GetCategory"

	c.Writer.Header().Set("content-type", "application/json")
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error get id category")
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.g.GetCategory(uint(id))
	if err != nil {
		h.log.WithField("op", op).Error("error get category")
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
//	@Param			req	body		RequestUpdateCategory					true	"данные для обновление категории"
//	@Success		200	{object}	api.SuccessResponse[uint]				"success"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некоректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректный данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Категория не найдена"
//
//	@Router			/category/update [put]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) UpdateCategory(c *gin.Context) {
	const op = "handlers.UpdateCategory"

	c.Writer.Header().Set("content-type", "application/json")
	var updateCategory RequestUpdateCategory

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	newCategory := domain.CategoryInput{
		Name:        updateCategory.Name,
		Limit:       updateCategory.Limit,
		Description: updateCategory.Description,
	}

	category, err := h.u.UpdateCategory(updateCategory.CategoryId, newCategory)
	if err != nil {
		h.log.WithField("op", op).Error("error update category")
		api.RegistrationError(c, err)
		return
	}
	api.ResponseOK(c, category)
}

// DeleteCategory godoc
//
//	@Summary		Удаление категории
//	@Description	Удаление категории конкретного пользователя
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			req	body		IDCategory								true	"userID для удаление категории"
//	@Success		200	{object}	api.SuccessResponse[string]				"success"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некоректные входные данные"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Категория не найдена"
//
//	@Router			/category/delete [delete]
//
//	@Security		jwtAuth
func (h *CategoryHandlers) DeleteCategory(c *gin.Context) {
	const op = "handlers.DeleteCategory"

	c.Writer.Header().Set("content-type", "application/json")
	param := c.Param("id")
	id, err := strconv.Atoi(param)

	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error get id")
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.d.DeleteCategory(uint(id))
	if err != nil {
		h.log.WithField("op", op).Error("error delete category")
		api.RegistrationError(c, err)
		return
	}
	api.ResponseOK(c, "user delete")
}
