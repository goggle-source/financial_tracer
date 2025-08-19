package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/financial_tracer/internal/domain"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DatabaseUserRepository interface {
	RegistrationUser(user *domain.User) (int, error)
	DeleteUser(email string, passwordHash []byte) error
	AuthenticationUser(email string, password []byte) (int, error)
}

type CreateUserServer struct {
	d DatabaseUserRepository
}

func CreateServer(d DatabaseUserRepository) *CreateUserServer {
	return &CreateUserServer{
		d: d,
	}
}

func (c *CreateUserServer) ServerRegistrationUser(name string, email string, password string) (int, error) {
	const op = "server.ServerCreateUser"
	if len(password) < 8 {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorPassword)
	}
	if !strings.Contains(email, "@") {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorEmail)
	}

	passwordHash, err := Hash(password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	user := domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	id, err := c.d.RegistrationUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrorDuplicated) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrorDuplicated)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c *CreateUserServer) ServerAuthenticationUser(email string, password string) (int, error) {
	const op = "server.ServerGetUser"
	if len(password) < 8 {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorPassword)
	}
	if !strings.Contains(email, "@") {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorEmail)
	}

	passwordHash, err := Hash(password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	id, err := c.d.AuthenticationUser(email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c *CreateUserServer) ServerDeleteUser(email string, password string) error {
	const op = "server.ServerDeleteUser"

	passwordHash, err := Hash(password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	err = c.d.DeleteUser(email, passwordHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
