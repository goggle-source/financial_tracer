package userHandlers

import (
	"net/http"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/handlers/api"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegistrationUserServic interface {
	RegistrationUser(us domain.RegisterUser) (jwttoken.ResponseJWTUser, error)
}

type AuthenticationUserServic interface {
	AuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error)
}

type DeleteUserServic interface {
	DeleteUser(us domain.DeleteUser) error
}

type HandlersUser struct {
	SecretKey string
	r         RegistrationUserServic
	a         AuthenticationUserServic
	d         DeleteUserServic
	log       *logrus.Logger
}

func CreateHandlersUser(secretKey string, r RegistrationUserServic,
	a AuthenticationUserServic,
	d DeleteUserServic, log *logrus.Logger) *HandlersUser {
	return &HandlersUser{
		d:         d,
		a:         a,
		r:         r,
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
//	@Param			user	body		UserRegistration						true	"Данные для регистрации пользователя"
//	@Success		200		{object}	api.SuccessResponse[ResponseJWTUser]	"Регистрация пользователя"
//
//	@Failure		400		{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500		{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400		{object}	api.ErrorResponse[string]				"Некорректные данные"
//
//	@Router			/registration/register [post]
func (h *HandlersUser) Registration(c *gin.Context) {
	const op = "handler.RegistrationUser"

	c.Writer.Header().Set("content-type", "application/json")
	var req UserRegistration

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	user := domain.RegisterUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := h.r.RegistrationUser(user)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error registration user")
		api.RegistrationError(c, err)
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
//	@Param			credentials	body		UserAuthentication						true	"Данные для авторизации пользователя"
//	@Success		200			{object}	api.SuccessResponse[ResponseJWTUser]	"Авторизация пользователя"
//
//	@Failure		400			{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500			{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400			{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404			{object}	api.ErrorResponse[string]				"Пользователь не найден"
//
//	@Router			/registration/login [post]
func (h *HandlersUser) Authentication(c *gin.Context) {
	const op = "handlers.GetUser"

	c.Writer.Header().Set("content-type", "application/json")
	var req UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	user := domain.AuthenticationUser{
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := h.a.AuthenticationUser(user)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error authentication user")
		api.RegistrationError(c, err)
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
//	@Param			req	body		UserDelete								true	"данные для удаление пользователя"
//	@Success		200	{object}	api.SuccessResponse[string]				"Удаление пользователя"
//
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Пользователь не найден"
//
//	@Router			/user/delete [post]
//
//	@Security		jwtAuth
func (h *HandlersUser) DeleteUser(c *gin.Context) {
	const op = "handlers.DeleteUser"

	c.Writer.Header().Set("content-type", "application/json")
	var req UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	users := domain.DeleteUser{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.d.DeleteUser(users)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error delete user")
		api.RegistrationError(c, err)
		return

	}

	api.ResponseOK(c, "user delete")
}

// GetAccessToken godoc
//
//	@Summary		Получение токена
//	@Description	Получение токена
//
//	@Tags			registration
//
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RefreshToken				true	"для получение access токена"
//	@Success		200	{object}	api.SuccessResponse[string]	"Получение access токена"
//
//	@Failure		400	{object}	api.ErrorResponse[string]	"Некорректные входные данные"
//
//	@Failure		500	{object}	api.ErrorResponse[string]	"Ошибка сервера"
//
//	@Router			/registration/get_access_token [post]
func (h *HandlersUser) GetAccessToken(c *gin.Context) {
	const op = "handlers.GetAccsessToken"

	c.Writer.Header().Set("content-type", "application/json")
	var req RefreshToken

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	i, name, err := jwttoken.CheckAccess(req.RefreshToken, h.SecretKey, h.log)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error check token")
		api.ResponseError(c, http.StatusBadRequest, "error valid token")
		return
	}

	tokens, err := jwttoken.JWTAccessToken(h.SecretKey, i, name)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	h.log.Info("create access token")
	api.ResponseOK(c, tokens)
}
