package category

import "errors"

var (
	ErrDatabase     = errors.New("error database")
	ErrNoFound      = errors.New("category is not found")
	ErrDuplicated   = errors.New("category is duplicated")
	ErrValidateType = errors.New("invalid type category")
)
