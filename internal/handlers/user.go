package handlers

import (
	"net/http"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/financial_tracer/internal/usercase/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandlersUser struct {
	users *user.CreateUserServer
	log   *logrus.Logger
}

func CreateHandlersUser(user *user.CreateUserServer, log *logrus.Logger) *HandlersUser {
	return &HandlersUser{
		users: user,
	}
}

func (h *HandlersUser) Post(c *gin.Context) {
	const op = "handler.PostUser"

	h.log.WithFields(logrus.Fields{
		"op": op,
	}).Info("post user")

	var user models.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorValidData,
		}).Error("Error request")
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrorValidData})
	}

	users, err := h.users.ServerRegistrationUser(user)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"op":  op,
			"err": domain.ErrorInternal,
		}).Error("error server")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error server"})
	}

	c.JSON(http.StatusOK, gin.H{"user": users})

}

func (h *HandlersUser) Get(c *gin.Context) {
	const op = "handlers.GetUser"

	h.log.WithFields(logrus.Fields{
		"op": op,
	}).Info("get user")

	var user models.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {

	}
}
