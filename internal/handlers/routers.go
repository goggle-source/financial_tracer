package handlers

import (
	_ "github.com/financial_tracer/docs"
	categoryHandlers "github.com/financial_tracer/internal/handlers/categories"
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
// @BasePath		/financial_tracker
func Router(users *userHandlers.HandlersUser, category *categoryHandlers.CategoryHandlers, log *logrus.Logger, sercretKey string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/financial_tracker")
	go api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	registration := api.Group("/registration")
	go registration.Use(middlewares.Logging(log))
	{

		go registration.POST("/register", users.Registration)
		go registration.POST("/login", users.Authentication)

		go registration.POST("/get_accsess_token", users.GetAccsessToken)
	}

	user := api.Group("/user")
	go user.Use(middlewares.Logging(log))
	go user.Use(middlewares.JWToken(sercretKey))
	{
		go user.POST("/delete", users.DeleteUser)
	}

	categories := api.Group("/category")
	go categories.Use(middlewares.Logging(log))
	go categories.Use(middlewares.JWToken(sercretKey))
	go categories.Use(middlewares.CORSMiddleware())
	{
		go categories.GET("/get_category", category.GetCategory)
		go categories.POST("/create_category", category.PostCategory)
		go categories.PUT("/update_category", category.UpdateCategory)
		go categories.DELETE("/delete_category", category.DeleteCategory)
	}

	return r
}
