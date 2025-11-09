package transaction

import (
	"errors"

	"github.com/financial_tracer/internal/infastructure/db/postgresql"
)

var (
	ErrNoFound  = errors.New("transaction is not found")
	ErrLimit    = errors.New("exceeded the limit")
	ErrDatabase = errors.New("error database")
)

func RegisterErrDatabase(err error) error {
	arr := map[error]error{
		postgresql.ErrorNotFound: ErrNoFound,
		postgresql.ErrorLimit:    ErrLimit,
	}

	value, ok := arr[err]
	if !ok {
		return ErrDatabase
	}

	return value
}
