package handlers

import (
	"github.com/financial_tracer/internal/handlers/middlewares"
	"github.com/gin-gonic/gin"
)

func Router(handlers *HandlersUser) *gin.Engine {
	r := gin.Default()

	registration := r.Group("/registration")
	registration.Use(middlewares.Logging(handlers.log))
	{
		registration.POST("/register", handlers.Registration)
		registration.POST("/login", handlers.Authentication)
	}

	user := r.Group("/user")
	user.Use(middlewares.Logging(handlers.log))
	user.Use(middlewares.JWToken(handlers.SecretKey))
	{
		user.POST("/delete", handlers.DeleteUser)
	}

	return r
}
