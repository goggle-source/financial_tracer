package transaction

import "errors"

var (
	ErrNoFound  = errors.New("transaction is not found")
	ErrLimit    = errors.New("exceeded the limit")
	ErrDatabase = errors.New("error database")
)
