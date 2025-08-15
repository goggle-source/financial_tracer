package domain

type User struct {
	RequestId    string
	Name         string
	Email        string
	PasswordHash []byte
}
