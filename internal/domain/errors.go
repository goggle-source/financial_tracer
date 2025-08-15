package domain

import "errors"

var (
	ErrorNotFound     = errors.New("not found")
	ErrorDuplicated   = errors.New("error enique")
	ErrorValidData    = errors.New("error request data")
	ErrorHashPassword = errors.New("error hash password")
	ErrorInternal     = errors.New("error server")
)
