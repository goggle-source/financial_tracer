package user

import "errors"

var (
	ErrDatabase   = errors.New("error database")
	ErrServic     = errors.New("servic error")
	ErrDuplicated = errors.New("the email has already been registered")
	ErrNoFound    = errors.New("user is not found")
)
