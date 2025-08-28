package user

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/go-playground/validator/v10"
)

type DatabaseUserRepository interface {
	RegistrationUser(user *domain.User) (int, string, error)
	DeleteUser(email string, password string) error
	AuthenticationUser(email string, password string) (int, string, error)
}

type UserValid struct {
	Valid func(error) []validator.ValidationErrors
}

type CreateUserServer struct {
	d DatabaseUserRepository
}

func CreateServer(d DatabaseUserRepository) *CreateUserServer {
	return &CreateUserServer{
		d: d,
	}
}

func (c *CreateUserServer) ServerRegistrationUser(us models.RegisterUser) (int, string, error) {
	const op = "server.ServerCreateUser"
	if err := validator.New().Struct(us); err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, ValidError(err))
	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
		PasswordHash: passwordHash,
	}

	id, name, err := c.d.RegistrationUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrorDuplicated) {
			return 0, "", fmt.Errorf("%s: %w", op, domain.ErrorDuplicated)
		}
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return id, name, nil
}

func (c *CreateUserServer) ServerAuthenticationUser(us models.User) (int, string, error) {
	const op = "server.ServerGetUser"
	if err := validator.New().Struct(us); err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, ValidError(err))
	}

	id, name, err := c.d.AuthenticationUser(us.Email, us.Password)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return 0, "", fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return id, name, nil
}

func (c *CreateUserServer) ServerDeleteUser(us models.User) error {
	const op = "server.ServerDeleteUser"
	if err := validator.New().Struct(us); err != nil {
		return fmt.Errorf("%s: %w", op, ValidError(err))
	}

	err := c.d.DeleteUser(us.Email, us.Password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func ValidError(err error) error {

	for _, err := range err.(validator.ValidationErrors) {

		switch err.Tag() {
		case "required":
			return domain.ErrorNotFound
		case "email":
			return domain.ErrorEmail
		case "min":
			return domain.ErrorSize
		case "max":
			return domain.ErrorSize
		default:
			return domain.ErrorValidData
		}
	}
	return nil
}
