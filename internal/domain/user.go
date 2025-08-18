package domain

type User struct {
	Name         string
	Email        string
	PasswordHash []byte
}
