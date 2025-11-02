package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Name	SuccessResponse
type SuccessResponse struct {
	Value any `json:"value,omitempty"`
}

func ResponseOK(c *gin.Context, value any) {
	c.Writer.Header().Set("content-type", "application/json")
	c.JSON(http.StatusOK, SuccessResponse{
		Value: value,
	})
}

// @Name	ErrorResponse
type ErrorResponse struct {
	MessageError any `json:"error"`
}

func ResponseError(c *gin.Context, status int, messageError any) {
	c.Writer.Header().Set("content-type", "application/json")
	c.JSON(status, ErrorResponse{
		MessageError: messageError,
	})
}

func ResponseUnauthorizedError(err string) ErrorResponse {
	return ErrorResponse{
		MessageError: err,
	}
}
