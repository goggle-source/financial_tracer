package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/financial_tracer/internal/servic/category"
	"github.com/financial_tracer/internal/servic/transaction"
	"github.com/financial_tracer/internal/servic/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errInfo struct {
	code    int
	message string
}

func RegistrationError(c *gin.Context, err error) {
	errMessage := ValidationError(err)
	if len(errMessage) != 0 {
		ResponseError(c, http.StatusBadRequest, errMessage)
		return
	}

	value := validateClientsErrors(err)
	fmt.Println(value)
	prov := errInfo{}
	if value == prov {
		ResponseError(c, http.StatusInternalServerError, "server error")
		return
	}

	ResponseError(c, value.code, value.message)

}

func ValidationError(err error) []map[string]string {
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

		return errMessage
	}

	return []map[string]string{}
}

func validateClientsErrors(err error) errInfo {

	arr := map[error]errInfo{
		transaction.ErrLimit: {
			code:    http.StatusBadRequest,
			message: "over the limit",
		},

		transaction.ErrNoFound: {
			code:    http.StatusNotFound,
			message: "transaction is not found",
		},

		category.ErrValidateType: {
			code:    http.StatusBadRequest,
			message: "param is not valid",
		},

		category.ErrDuplicated: {
			code:    http.StatusBadRequest,
			message: "this category already exists",
		},

		category.ErrNoFound: {
			code:    http.StatusNotFound,
			message: "category is not fuond",
		},

		user.ErrDuplicated: {
			code:    http.StatusBadRequest,
			message: "this user already exists",
		},

		user.ErrNoFound: {
			code:    http.StatusNotFound,
			message: "user is not found",
		},

		user.ErrServic: {
			code:    http.StatusInternalServerError,
			message: "server error",
		},

		user.ErrDatabase: {
			code:    http.StatusInternalServerError,
			message: "server error",
		},

		category.ErrDatabase: {
			code:    http.StatusInternalServerError,
			message: "server error",
		},

		transaction.ErrDatabase: {
			code:    http.StatusInternalServerError,
			message: "server error",
		},
	}

	value, ok := arr[err]
	if !ok {
		return errInfo{}
	}

	return value
}
