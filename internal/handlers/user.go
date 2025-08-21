package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/financial_tracer/internal/servic/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type ResponseJSONUser struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AccsessToken string `json:"access_token"`
}

type Claims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
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

//  Registration godoc
//	@Summary		Регистрция пользователя
//	@Description	Регистрация пользователя
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			req	body		models.RegisterUser	true	"RegisterUser"
//	@Success		200	{object}	ResponseJSONUser
//	@Router			/registration/register [post]

func (h *HandlersUser) Registration(c *gin.Context) {
	const op = "handler.PostUser"

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorValidData,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	user := models.RegisterUser{
		Name: req.Name,
		User: models.User{Email: req.Email, Password: req.Password},
	}

	userId, err := h.users.ServerRegistrationUser(user)
	if err != nil {
		RegistrationError(c, op, err)

		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error server")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	tokens, err := PostJWT(c, h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("registration user")
	c.JSON(http.StatusOK, tokens)

}

//  Authentication godoc
//	@Summary		Аутентификация пользователя
//	@Description	Аутентификация пользователя
//	@Accept			json
//	@Produce		json
//	@Param			req	body		models.User	true	"User"
//	@Success		200	{object}	ResponseJSONUser
//	@Router			/registration/login [post]

func (h *HandlersUser) Authentication(c *gin.Context) {
	const op = "handlers.GetUser"

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorValidData,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	userId, err := h.users.ServerAuthenticationUser(user)
	if err != nil {

		RegistrationError(c, op, err)

		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error authentication user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	tokens, err := PostJWT(c, h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("registration user")
	c.JSON(http.StatusOK, tokens)
}

//  DeleteUser godoc
//	@Summary		Удаление пользователя
//	@Description	Удаление пользователя
//	@Accept			json
//	@Produce		json
//	@Param			req	body		models.User	true	"User"
//	@Success		200	{object}	string
//	@Router			/user/delete [post]
func (h *HandlersUser) DeleteUser(c *gin.Context) {
	const op = "handlers.DeleteUser"

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorValidData,
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
				"err": domain.ErrorNotFound,
			}).Error("error not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "error not found"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

//  GetAccsessToken godoc
//	@summary		Получение токена
//	@description	Получение токена
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
			"err": domain.ErrorValidData,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrorInternal)
		}

		return []byte(h.SecretKey), nil
	})

	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error parse token",
		}).Error("error parse token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	if !token.Valid {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error valid token",
		}).Error("error valid token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	tokenClaims, ok := token.Claims.(*Claims)
	if !ok {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error claims",
		}).Error("error claims")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	if tokenClaims.ExpiresAt.Unix() < time.Now().Unix() {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error expired token",
		}).Error("error expired token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	if tokenClaims.ID == "" {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error id",
		}).Error("error id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	i, err := strconv.Atoi(tokenClaims.ID)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error id",
		}).Error("error id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error token"})
		return
	}

	tokens, err := JWTAccsessToken(h.SecretKey, i)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error create tokens",
		}).Error("error create tokens")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("create accsess token")
	c.JSON(http.StatusOK, gin.H{"accsess_token": tokens})
}
