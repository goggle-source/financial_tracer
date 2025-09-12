package userHandlers

import (
	"net/http"

	"github.com/financial_tracer/internal/handlers/api"
	"github.com/financial_tracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServicUserer interface {
	ServerRegistrationUser(us models.RegisterUser) (uint, string, error)
	ServerAuthenticationUser(us models.AuthenticationUser) (uint, string, error)
	ServerDeleteUser(us models.DeleteUser) error
}

type HandlersUser struct {
	SecretKey string
	users     ServicUserer
	log       *logrus.Logger
}

func CreateHandlersUser(secretKey string, user ServicUserer, log *logrus.Logger) *HandlersUser {
	return &HandlersUser{
		users:     user,
		log:       log,
		SecretKey: secretKey,
	}
}

// RegistrationUser godoc
//
//	@Summary		Регистрация пользователя
//	@Description	Создание нового пользователя
//	@Tags			registration
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserRegistration	true	"Данные для регистрации пользователя"
//	@Success		200		{object}	api.SuccessResponse[ResponseJWTUser] "Регистрация пользователя"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректные данные"
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/registration/register [post]
func (h *HandlersUser) Registration(c *gin.Context) {
	const op = "handler.RegistrationUser"

	var req UserRegistration

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	user := models.RegisterUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	userId, userName, err := h.users.ServerRegistrationUser(user)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error registration user")
		api.RegistrationError(c, op, err)
		return
	}

	tokens, err := PostJWT(c, h.SecretKey, userId, userName)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	api.ResponseOK(c, tokens)

}

// AuthenticationUser godoc
//
//	@Summary		Аутентификация пользователя
//	@Description	Вход пользователя в систему
//	@Tags			registration
//
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		UserAuthentication	true	"Данные для авторизации пользователя"
//	@Success		200			{object}	api.SuccessResponse[ResponseJWTUser] "Авторизация пользователя"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректные данные"
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/registration/login [post]
func (h *HandlersUser) Authentication(c *gin.Context) {
	const op = "handlers.GetUser"

	var req UserAuthentication

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	user := models.AuthenticationUser{
		Email:    req.Email,
		Password: req.Password,
	}

	userId, userName, err := h.users.ServerAuthenticationUser(user)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error authentication user")
		api.RegistrationError(c, op, err)
		return

	}

	tokens, err := PostJWT(c, h.SecretKey, userId, userName)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	api.ResponseOK(c, tokens)
}

// DeleteUser godoc
//
//	@Summary		Удаление пользователя
//	@Description	Удаление пользователя
//
//	@Tags			User
//
//	@Accept			json
//	@Produce		json
//	@Param			req	body		UserDelete	true	"данные для удаление пользователя"
//	@Success		200	{object}	api.SuccessResponse[string] "Удаление пользователя"
//
// @Failure 400 {object} api.ErrorResponse[[]map[string]string] "Некоректные данные"
// @Failure 500 {object} api.ErrorResponse[string] "Ошибка сервера"
// @Failure 400 {object} api.ErrorResponse[string] "Некоректные данные"
//
//	@Router			/user/delete [post]
func (h *HandlersUser) DeleteUser(c *gin.Context) {
	const op = "handlers.DeleteUser"

	var req UserDelete

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	users := models.DeleteUser{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.users.ServerDeleteUser(users)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error delete user")
		api.RegistrationError(c, op, err)
		return

	}

	api.ResponseOK(c, "user delete")
}

// GetAccsessToken godoc
//
//	@Summary		Получение токена
//	@Description	Получение токена
//
//	@Tags			registration
//
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RefreshToken	true	"для получение accsess токена"
//	@Success		200	{object}	api.SuccessResponse[string] "Получение access токена"
//	@Router			/registration/get_accsess_token [post]
func (h *HandlersUser) GetAccsessToken(c *gin.Context) {
	const op = "handlers.GetAccsessToken"

	var req RefreshToken

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	i, name, err := CheckAccess(req.RefreshToken, h.SecretKey, h.log)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error check token")
		api.ResponseError(c, http.StatusBadRequest, "error valid token")
		return
	}

	tokens, err := JWTAccessToken(h.SecretKey, i, name)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	h.log.Info("create accsess token")
	api.ResponseOK(c, tokens)
}
