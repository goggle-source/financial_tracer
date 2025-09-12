package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/financial_tracer/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegistrationError(c *gin.Context, op string, err error) {
	var validError validator.ValidationErrors
	if errors.As(err, &validError) {
		var errMessage []map[string]string
		for _, errs := range validError {
			detail := map[string]string{
				"field":   errs.Field(),
				"message": fmt.Sprintf(" Field validation for %s, field on the tag: %s", errs.Field(), errs.Tag()),
			}

			errMessage = append(errMessage, detail)
		}

		ResponseError(c, http.StatusBadRequest, errMessage)
		return
	}

	if errMsg := ClientError(err); errMsg != "" {
		ResponseError(c, http.StatusBadRequest, errMsg)
		return
	}

	ResponseError(c, http.StatusInternalServerError, "error server")

}

func ClientError(err error) string {
	clientErrors := map[error]string{
		domain.ErrorNotFound:   "not found",
		domain.ErrorDuplicated: "duplicated unique",
		domain.ErrorValidData:  "error request",
	}
	value, ok := clientErrors[err]
	if !ok {
		return ""
	}
	return value
}
