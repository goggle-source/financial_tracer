package domain

type User struct {
	Name         string
	Email        string
	PasswordHash []byte
}

type Category struct {
	Name        string
	Limit       int
	Description string
}
