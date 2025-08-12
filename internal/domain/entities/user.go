package entities

import (
	"errors"
	"strings"
)

var (
	ErrorUserExists    = errors.New("user exists")
	ErrorValidEmail    = errors.New("not valid email")
	ErrorValidPassword = errors.New("not valid password")
	ErrorValidName     = errors.New("not valid name")
)

type User struct {
	Id           uint
	RequestId    string
	Name         string
	Email        string
	PasswordHash [32]byte
}

func NewUser(requestId string, name string, email string, passwordHash [32]byte) (*User, error) {
	if ValidEmail(email) {
		return nil, ErrorValidEmail
	}

	if ValidName(name) {
		return nil, ErrorValidName
	}

	return &User{
		RequestId:    requestId,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func ValidEmail(email string) bool {
	if email == "" && !strings.Contains(email, "@") {
		return false
	}
	return true
}

func ValidName(name string) bool {
	if len(name) >= 50 && name == "" {
		return false
	}
	return true
}
