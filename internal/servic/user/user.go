package user

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/go-playground/validator/v10"
)

type DatabaseUserRepository interface {
	RegistrationUser(user *domain.User) (int, error)
	DeleteUser(email string, passwordHash []byte) error
	AuthenticationUser(email string, password []byte) (int, error)
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

func (c *CreateUserServer) ServerRegistrationUser(us models.RegisterUser) (int, error) {
	const op = "server.ServerCreateUser"
	if err := validator.New().Struct(us); err != nil {

	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
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

func (c *CreateUserServer) ServerAuthenticationUser(us models.User) (int, error) {
	const op = "server.ServerGetUser"
	if err := validator.New().Struct(us); err != nil {
		return 0, fmt.Errorf("%s: %w", op, ValidError(err))
	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	id, err := c.d.AuthenticationUser(us.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c *CreateUserServer) ServerDeleteUser(us models.User) error {
	const op = "server.ServerDeleteUser"
	if err := validator.New().Struct(us); err != nil {
		return fmt.Errorf("%s: %w", op, ValidError(err))
	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	err = c.d.DeleteUser(us.Email, passwordHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func ValidError(err error) error {

	for _, err := range err.(validator.ValidationErrors) {

		switch err.Tag() {
		case "required":
			return fmt.Errorf("%s is required", err.Field())
		case "email":
			return fmt.Errorf("%s is not valid email", err.Field())
		case "min":
			return fmt.Errorf("%s must be at least %s characters", err.Field(), err.Param())
		case "max":
			return fmt.Errorf("%s must be at most %s characters", err.Field(), err.Param())
		default:
			return fmt.Errorf("%s is not valid", err.Field())
		}
	}
	return nil
}
