package userHandlers

import (
	"errors"
	"net/http"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/financial_tracer/internal/servic/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserRegistration represents registration user request
type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserAuthentication represents authentication user request
type UserAuthentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserDelete represents delete user request
type UserDelete struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResponseJSONUser represents response user response
type ResponseJSONUser struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AccsessToken string `json:"access_token"`
}

type HandlersUser struct {
	SecretKey string
	users     *user.CreateUserServer
	log       *logrus.Logger
}

func CreateHandlersUser(secretKey string, user *user.CreateUserServer, log *logrus.Logger) *HandlersUser {
	return &HandlersUser{
		users:     user,
		log:       log,
		SecretKey: secretKey,
	}
}

// Registration godoc
//
//	@Summary		Регистрация пользователя
//	@Description	Создание нового пользователя
//	@Tags			registration
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UserRegistration	true	"Данные пользователя"
//	@Success		201		{object}	ResponseJSONUser
//	@Router			/registration/register [post]
func (h *HandlersUser) Registration(c *gin.Context) {
	const op = "handler.PostUser"

	var req UserRegistration

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	user := models.RegisterUser{
		Name: req.Name,
		User: models.User{Email: req.Email, Password: req.Password},
	}

	userId, userName, err := h.users.ServerRegistrationUser(user)
	if err != nil {
		RegistrationError(c, op, err)

		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error server")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	tokens, err := PostJWT(c, h.SecretKey, userId, userName)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("registration user")
	c.JSON(http.StatusOK, tokens)

}

// Authentication godoc
//
//	@Summary		Аутентификация пользователя
//	@Description	Вход пользователя в систему
//	@Tags			registration
//
// @Accept			json
// @Produce		json
// @Param			credentials	body		UserAuthentication	true	"Учетные данные"
// @Success		200			{object}	ResponseJSONUser
// @Router			/registration/login [post]
func (h *HandlersUser) Authentication(c *gin.Context) {
	const op = "handlers.GetUser"

	var req UserAuthentication

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	userId, userName, err := h.users.ServerAuthenticationUser(user)
	if err != nil {

		if RegistrationError(c, op, err) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": err,
			}).Error("error authentication user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
			return
		}
		return

	}

	tokens, err := PostJWT(c, h.SecretKey, userId, userName)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("registration user")
	c.JSON(http.StatusOK, tokens)
}

// DeleteUser godoc
//
//	@Summary		Удаление пользователя
//	@Description	Удаление пользователя
//
// @Tags deleteUser
//
//	@Accept			json
//	@Produce		json
//	@Param			req	body		UserDelete	true	"User"
//	@Success		200	{object}	string
//	@Router			/user/delete [post]
func (h *HandlersUser) DeleteUser(c *gin.Context) {
	const op = "handlers.DeleteUser"

	var req UserDelete

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	users := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.users.ServerDeleteUser(users)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": err,
			}).Error("error not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "error not found"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// GetAccsessToken godoc
//
//	@summary		Получение токена
//	@description	Получение токена
//
// @Tags registration
//
//	@accept			json
//	@produce		json
//	@param			req	body		models.User	true	"User"
//	@success		200	{object}	string
//	@router			/registration/get_accsess_token [post]
func (h *HandlersUser) GetAccsessToken(c *gin.Context) {
	const op = "handlers.GetAccsessToken"

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	i, name, err := CheckAccess(req.RefreshToken, h.SecretKey, h.log)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":    op,
			"error": err,
		}).Error("error check token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is not valid"})
	}

	tokens, err := JWTAccessToken(h.SecretKey, i, name)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("create accsess token")
	c.JSON(http.StatusOK, gin.H{"accsess_token": tokens})
}
