package categoryHandlers

import (
	"net/http"

	"github.com/financial_tracer/internal/handlers/api"
	"github.com/financial_tracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServicCategoryer interface {
	Create(idUser uint, category models.Category) (uint, error)
	Get(idUser uint, idCategory uint) (models.Category, error)
	Update(idUser uint, idCategory uint, newCategory models.Category) (uint, error)
	Delete(idUser uint, idCategory uint) error
}

type CategoryHandlers struct {
	categor ServicCategoryer
	log     *logrus.Logger
}

func CreateHandlersCategory(categ ServicCategoryer, log *logrus.Logger) *CategoryHandlers {
	return &CategoryHandlers{
		categor: categ,
		log:     log,
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
// @Success 200 {object} api.SuccessResponse[uint] "Успешное создание категории"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректено введены данные"
//
// @Failure 501 {object} api.ErrorResponse[string] "Ошибка сервера"
//
// @Router			/category/create_category [post]
func (h *CategoryHandlers) PostCategory(c *gin.Context) {
	const op = "handlers.PostCategory"

	newCategory := RequestCreateCategory{}
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

	cat := models.Category{
		Name:        newCategory.Name,
		Limit:       newCategory.Limit,
		Description: newCategory.Description,
	}

	idCategory, err := h.categor.Create(idUser.(uint), cat)
	if err != nil {
		h.log.WithField("op", op).Error("error create category")
		api.RegistrationError(c, op, err)
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
//	@Param			req			body		IDCategory	true	"userID для получение категории"
//	@Success		200 {object}	api.SuccessResponse[models.Category]	"success"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректные данные"
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/category/get_category [get]
func (h *CategoryHandlers) GetCategory(c *gin.Context) {
	const op = "handlers.GetCateogry"

	id := IDCategory{}
	if err := c.ShouldBindJSON(&id); err != nil {
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

	category, err := h.categor.Get(idUser.(uint), id.CategoryId)
	if err != nil {
		h.log.WithField("op", op).Error("error get category")
		api.RegistrationError(c, op, err)
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
//	@Success		200	{object}	api.SuccessResponse[uint]				"success"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректные данные"
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/category/update_category [put]
func (h *CategoryHandlers) UpdateCategory(c *gin.Context) {
	const op = "handlers.UpdateCategory"

	var updateCategory RequestUpdateCategory

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	newCategory := models.Category{
		Name:        updateCategory.Name,
		Limit:       updateCategory.Limit,
		Description: updateCategory.Description,
	}

	idUser, ok := c.Get("userID")
	if !ok {
		h.log.WithField("op", op).Error("error get userID")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	id, err := h.categor.Update(idUser.(uint), updateCategory.CategoryId, newCategory)
	if err != nil {
		h.log.WithField("op", op).Error("error update category")
		api.RegistrationError(c, op, err)
		return
	}
	api.ResponseOK(c, id)
}

// DeleteCategory godoc
//
//	@Summary		Удаление категории
//	@Description	Удаление категории конкретного пользователя
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			req	body		IDCategory	true	"userID для удаление категории"
//	@Success		200	{object}	api.SuccessResponse[string]	"success"
//
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/category/delete_category [delete]
func (h *CategoryHandlers) DeleteCategory(c *gin.Context) {
	const op = "handlers.DeleteCategory"

	categoryId := IDCategory{}
	if err := c.ShouldBindJSON(&categoryId); err != nil {
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

	err := h.categor.Delete(idUser.(uint), categoryId.CategoryId)
	if err != nil {
		h.log.WithField("op", op).Error("error delete category")
		api.RegistrationError(c, op, err)
		return
	}
	api.ResponseOK(c, "user delete")
}
