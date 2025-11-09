package category

import (
	"errors"

	"github.com/financial_tracer/internal/infastructure/db/postgresql"
)

var (
	ErrDatabase     = errors.New("error database")
	ErrNoFound      = errors.New("category is not found")
	ErrDuplicated   = errors.New("category is duplicated")
	ErrValidateType = errors.New("invalid type category")
)

func RegsiterErrorDatabase(err error) error {
	arr := map[error]error{
		postgresql.ErrorDuplicated: ErrDuplicated,
		postgresql.ErrorNotFound:   ErrNoFound,
		postgresql.ErrorLimit:      ErrValidateType,
	}

	value, ok := arr[err]
	if !ok {
		return ErrDatabase
	}

	return value
}
