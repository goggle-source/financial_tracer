package handlers

import (
	"github.com/financial_tracer/internal/handlers/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Финансый Трекер
// @version		1.0
// @description	API для работы с финансовым трекером
// @license.name	MIT
// @contact.name	Ярослав
// @contact.url	https://github.com/goggle-source
// @contact.email	asssv0423348@gmail.com
// @host			localhost:8080
// @BasePath		/
func Router(handlers *HandlersUser) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registration := r.Group("/registration")
	registration.Use(middlewares.Logging(handlers.log))
	{

		registration.POST("/register", handlers.Registration)
		registration.POST("/login", handlers.Authentication)

		registration.POST("/get_accsess_token", handlers.GetAccsessToken)
	}

	user := r.Group("/user")
	user.Use(middlewares.Logging(handlers.log))
	user.Use(middlewares.JWToken(handlers.SecretKey))
	{
		user.POST("/delete", handlers.DeleteUser)
	}

	return r
}
