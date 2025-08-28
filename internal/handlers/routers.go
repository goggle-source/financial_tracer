package handlers

import (
	"github.com/financial_tracer/internal/handlers/middlewares"
	userHandlers "github.com/financial_tracer/internal/handlers/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Финансовый Трекер
// @version		1.0
// @description	API для работы с финансовым трекером
// @license.name	MIT
// @contact.name	Ярослав
// @contact.url	https://github.com/goggle-source
// @contact.email	asssv0423348@gmail.com
// @host			localhost:8080
// @BasePath		/
func Router(users *userHandlers.HandlersUser, log *logrus.Logger) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registration := r.Group("/registration")
	registration.Use(middlewares.Logging(log))
	{

		registration.POST("/register", users.Registration)
		registration.POST("/login", users.Authentication)

		registration.POST("/get_accsess_token", users.GetAccsessToken)
	}

	user := r.Group("/user")
	user.Use(middlewares.Logging(log))
	user.Use(middlewares.JWToken(users.SecretKey))
	{
		user.POST("/delete", users.DeleteUser)
	}

	return r
}
