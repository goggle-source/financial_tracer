package handlers

import "github.com/gin-gonic/gin"

func Router(handlers *HandlersUser) *gin.Engine {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", handlers.Registration)
		user.POST("/login", handlers.Authentication)
	}

	return r
}
