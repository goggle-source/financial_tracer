package domain

import "errors"

var (
	ErrorNotFound     = errors.New("not found")
	ErrorDuplicated   = errors.New("duplicated user")
	ErrorValidData    = errors.New("invalid request")
	ErrorHashPassword = errors.New("error hash password")
	ErrorInternal     = errors.New("error server")
)
