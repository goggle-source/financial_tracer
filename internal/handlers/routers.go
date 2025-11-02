package handlers

import (
	_ "github.com/financial_tracer/docs"
	categoryHandlers "github.com/financial_tracer/internal/handlers/categories"
	"github.com/financial_tracer/internal/handlers/middlewares"
	transactionHandlers "github.com/financial_tracer/internal/handlers/transaction"
	userHandlers "github.com/financial_tracer/internal/handlers/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						Финансовый Трекер
// @version					1.0
// @description				API для работы с финансовым трекером
// @license.name				MIT
// @contact.name				Ярослав
// @contact.url				https://github.com/goggle-source
// @contact.email				asssv0423348@gmail.com
// @host						localhost:8080
// @BasePath					/financial_tracker
//
// @securityDefinitions.apiKey	jwtAuth
// @in							header
// @name						Authorization
// @description				type "Bearer" после пробел и jwt token, пример: "Bearer zpdgjeawzgp0398tuP29R0J20THVTP9235BHRNr312r346as2..."
func Router(users *userHandlers.HandlersUser, category *categoryHandlers.CategoryHandlers, log *logrus.Logger, tran *transactionHandlers.TransactionHandlers, secretKey string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/financial_tracker")
	go api.Use(middlewares.Logging(log))
	go api.Use(middlewares.CORSMiddleware())

	go api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	registration := api.Group("/registration")
	{

		go registration.POST("/register", users.Registration)
		go registration.POST("/login", users.Authentication)

		go registration.POST("/get_access_token", users.GetAccessToken)
	}

	user := api.Group("/user")
	go user.Use(middlewares.Logging(log))
	go user.Use(middlewares.JWToken(secretKey, log))
	{
		go user.POST("/delete", users.DeleteUser)
	}

	categories := api.Group("/category")
	go categories.Use(middlewares.JWToken(secretKey, log))
	{
		go categories.GET("/get/{id}", category.GetCategory)
		go categories.GET("/get/type/{id}", category.CategoryType)
		go categories.POST("/create", category.PostCategory)
		go categories.PUT("/update", category.UpdateCategory)
		go categories.DELETE("/delete{id}", category.DeleteCategory)
	}
	transaction := api.Group("/transaction")
	go transaction.Use(middlewares.JWToken(secretKey, log))
	{
		go transaction.POST("/create", tran.PostTransaction)
		go transaction.GET("/get/{id}", tran.GetTransaction)
		go transaction.PUT("/update", tran.UpdateTransaction)
		go transaction.DELETE("/delete/{id}", tran.DeleteTransaction)
	}

	return r
}
