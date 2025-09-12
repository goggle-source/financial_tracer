package user

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"github.com/go-playground/validator/v10"
)

type DatabaseUserRepository interface {
	RegistrationUser(user *domain.User) (uint, string, error)
	DeleteUser(email string, password string) error
	AuthenticationUser(email string, password string) (uint, string, error)
}

type UserValid struct {
	Valid func(error) []validator.ValidationErrors
}

type UserServer struct {
	d DatabaseUserRepository
}

func CreateUserServer(d DatabaseUserRepository) *UserServer {
	return &UserServer{
		d: d,
	}
}

func (c *UserServer) ServerRegistrationUser(us models.RegisterUser) (uint, string, error) {
	if err := validator.New().Struct(us); err != nil {
		return 0, "", err
	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return 0, "", err
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
		PasswordHash: passwordHash,
	}

	id, name, err := c.d.RegistrationUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrorDuplicated) {
			return 0, "", err
		}
		return 0, "", err
	}

	return id, name, nil
}

func (c *UserServer) ServerAuthenticationUser(us models.AuthenticationUser) (uint, string, error) {

	if err := validator.New().Struct(us); err != nil {

		return 0, "", err
	}

	id, name, err := c.d.AuthenticationUser(us.Email, us.Password)
	if err != nil {
		return 0, "", err
	}

	return id, name, nil
}

func (c *UserServer) ServerDeleteUser(us models.DeleteUser) error {
	if err := validator.New().Struct(us); err != nil {

		return err
	}

	err := c.d.DeleteUser(us.Email, us.Password)
	if err != nil {
		return err
	}
	return nil
}
