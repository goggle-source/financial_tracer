package postgresql

import "errors"

var (
	ErrorNotFound   = errors.New("not found")
	ErrorDuplicated = errors.New("duplicated unique")
	ErrorLimit      = errors.New("error limit transaction")
)
