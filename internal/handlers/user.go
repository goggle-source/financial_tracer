package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/servic/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

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

	userId, err := h.users.ServerRegistrationUser(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrorDuplicated) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorDuplicated,
			}).Error("error duplicated email")
			c.JSON(http.StatusConflict, gin.H{"error": "error duplicated email"})
			return
		}

		if errors.Is(err, domain.ErrorEmail) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorEmail,
			}).Error("error email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error email"})
			return
		}

		if errors.Is(err, domain.ErrorPassword) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorPassword,
			}).Error("error password")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error password"})
			return
		}

		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error server")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	t, err := JWTAccsessToken(h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create accsess token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	rt, err := JWTRefreshToken(h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("registration user")
	c.JSON(http.StatusOK, ResponseJSONUser{RefreshToken: rt, AccsessToken: t})

}

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

	userId, err := h.users.ServerAuthenticationUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorNotFound,
			}).Error("error not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "error not found"})
		}
		if errors.Is(err, domain.ErrorPassword) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorPassword,
			}).Error("error password")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error password"})
		}
		if errors.Is(err, domain.ErrorEmail) {
			h.log.WithFields(logrus.Fields{
				"op":  op,
				"err": domain.ErrorEmail,
			}).Error("error email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error email"})
			return
		}

		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error authentication user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	t, err := JWTAccsessToken(h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create accsess token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	rt, err := JWTRefreshToken(h.SecretKey, userId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error create refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
		return
	}

	h.log.Info("authentication user")
	c.JSON(http.StatusOK, ResponseJSONUser{RefreshToken: rt, AccsessToken: t})
}

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

	err := h.users.ServerDeleteUser(req.Email, req.Password)
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

func JWTAccsessToken(secretKey string, id int) (string, error) {
	const op = "handlers.JWTAccsessToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func JWTRefreshToken(secretKey string, id int) (string, error) {
	const op = "handlers.JWTAccsessToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 148).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}
