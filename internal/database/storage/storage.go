package storage

import "errors"

var (
	ErrorNotFound          = errors.New("not found")
	ErrorUserExists        = errors.New("user exists")
	ErrorCategoryExists    = errors.New("category exists")
	ErrorTransactionExists = errors.New("transaction exists")
)
