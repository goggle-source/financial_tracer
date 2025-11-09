package user

import (
	"errors"

	"github.com/financial_tracer/internal/infastructure/db/postgresql"
)

var (
	ErrDatabase   = errors.New("error database")
	ErrServic     = errors.New("servic error")
	ErrDuplicated = errors.New("the email has already been registered")
	ErrNoFound    = errors.New("user is not found")
)

func RegisterErrDatabase(err error) error {
	arr := map[error]error{
		postgresql.ErrorDuplicated: ErrDuplicated,
		postgresql.ErrorNotFound:   ErrNoFound,
	}

	value, ok := arr[err]
	if !ok {
		return ErrDatabase
	}

	return value
}
