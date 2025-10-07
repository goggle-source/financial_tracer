package postgresql

import "errors"

var (
	ErrorNotFound     = errors.New("not found")
	ErrorDuplicated   = errors.New("duplicated unique")
	ErrorValidData    = errors.New("invalid request")
	ErrorHashPassword = errors.New("error hash password")
	ErrorInternal     = errors.New("error server")
	ErrorLimit        = errors.New("error limit transaction")
)
