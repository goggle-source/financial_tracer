package handlers

import (
	"github.com/financial_tracer/docs"
	categoryHandlers "github.com/financial_tracer/internal/handlers/categories"
	"github.com/financial_tracer/internal/handlers/middlewares"
	transactionHandlers "github.com/financial_tracer/internal/handlers/transaction"
	userHandlers "github.com/financial_tracer/internal/handlers/user"
	"github.com/gin-contrib/pprof"
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
//
// @securityDefinitions.apiKey	jwtAuth
// @in							header
// @name						Authorization
// @description				type "Bearer" после пробел и jwt token, пример: "Bearer zpdgjeawzgp0398tuP29R0J20THVTP9235BHRNr312r346as2..."
func Router(users *userHandlers.HandlersUser, category *categoryHandlers.CategoryHandlers, log *logrus.Logger, tran *transactionHandlers.TransactionHandlers, secretKey string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/financial_tracker")
	api.Use(middlewares.Logging(log))
	api.Use(middlewares.CORSMiddleware())

	registration := api.Group("/registration")
	{

		registration.POST("/register", users.Registration)
		registration.POST("/login", users.Authentication)

		registration.POST("/access_token", users.GetAccessToken)
	}

	user := api.Group("/user")
	user.Use(middlewares.Logging(log))
	user.Use(middlewares.JWToken(secretKey, log))
	{
		user.DELETE("/", users.DeleteUser)
	}

	categories := api.Group("/category")
	categories.Use(middlewares.JWToken(secretKey, log))
	{
		categories.GET("/:id", category.GetCategory)
		categories.GET("/type/:type", category.CategoryType)
		categories.POST("/", category.PostCategory)
		categories.PUT("/", category.UpdateCategory)
		categories.DELETE("/:id", category.DeleteCategory)
	}

	transaction := api.Group("/transaction")
	transaction.Use(middlewares.JWToken(secretKey, log))
	{
		transaction.POST("/", tran.PostTransaction)
		transaction.GET("/:id", tran.GetTransaction)
		transaction.PUT("/", tran.UpdateTransaction)
		transaction.DELETE("/:id", tran.DeleteTransaction)
	}

	docs.SwaggerInfo.BasePath = "/financial_tracker"
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	pprof.Register(api, "/debug/pprof")

	return r
}
