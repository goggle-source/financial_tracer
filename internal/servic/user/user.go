package user

import (
	"fmt"

	"github.com/financial_tracer/internal/domain"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type DatabaseUserRepository interface {
	RegistrationUser(user domain.User) (uint, string, error)
	DeleteUser(email string, password string) error
	AuthenticationUser(email string, password string) (uint, string, error)
}

type UserValid struct {
	Valid func(error) []validator.ValidationErrors
}

type UserServer struct {
	log       *logrus.Logger
	d         DatabaseUserRepository
	SecretKey string
}

func CreateUserServer(d DatabaseUserRepository, sk string, log *logrus.Logger) *UserServer {
	return &UserServer{
		log:       log,
		d:         d,
		SecretKey: sk,
	}
}

func (c *UserServer) ServerRegistrationUser(us domain.RegisterUser) (jwttoken.ResponseJWTUser, error) {
	const op = "user.ServerRegistrationUser"

	if err := validator.New().Struct(us); err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	passwordHash, err := Hash(us.Password)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field hash password: %w", op, err)
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
		PasswordHash: passwordHash,
	}

	id, name, err := c.d.RegistrationUser(user)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field registration user: %w", op, err)
	}

	tokens, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field create JWT token: %w", op, err)
	}

	return tokens, nil
}

func (c *UserServer) ServerAuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	const op = "user.ServerAuthenticationUser"

	if err := validator.New().Struct(us); err != nil {

		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	id, name, err := c.d.AuthenticationUser(us.Email, us.Password)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field authentication user: %w", op, err)
	}

	token, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field create JWT token: %w", op, err)
	}

	return token, nil
}

func (c *UserServer) ServerDeleteUser(us domain.DeleteUser) error {
	if err := validator.New().Struct(us); err != nil {

		return fmt.Errorf("invalid validate: %w", err)
	}

	err := c.d.DeleteUser(us.Email, us.Password)
	if err != nil {
		return fmt.Errorf("field delete user: %w", err)
	}
	return nil
}
