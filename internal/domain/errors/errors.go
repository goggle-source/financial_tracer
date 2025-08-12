package storage

import "errors"

var (
	ErrorNotFound          = errors.New("not found")
	ErrorDuplicated        = errors.New("error enique")
	ErrorCategoryExists    = errors.New("category exists")
	ErrorTransactionExists = errors.New("transaction exists")
)
