package user

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/financial_tracer/internal/lib/hashPassword"
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

	log := c.log.WithField("op", op)

	log.Info("registration user")

	if err := validator.New().Struct(us); err != nil {
		log.WithFields(logrus.Fields{"err": err}).Error("invalid validate")

		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	passwordHash, err := hashPassword.Hash(us.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("field hash password")

		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field hash password: %w", op, ErrServic)
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
		PasswordHash: passwordHash,
	}

	id, name, err := c.d.RegistrationUser(user)
	if err != nil {
		if errors.Is(err, postgresql.ErrorDuplicated) {
			log.WithField("err", err).Error("field duplicated")

			return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field duplicated: %w", op, ErrDuplicated)
		}
		log.WithField("err", err).Error("field registration user")
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field registration user: %w", op, ErrDatabase)
	}

	tokens, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		log.WithField("err", err).Error("field create JWT token")
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field create JWT token: %w", op, ErrServic)
	}

	log.Info("success registration user")

	return tokens, nil
}

func (c *UserServer) ServerAuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	const op = "user.ServerAuthenticationUser"

	log := c.log.WithField("op", op)

	log.Info("start authentication user")

	if err := validator.New().Struct(us); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	id, name, err := c.d.AuthenticationUser(us.Email, us.Password)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("not found")

			return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("field authentication user")
		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field authentication user: %w", op, ErrDatabase)
	}

	token, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		log.WithField("err", err).Error("field create JWT token")

		return jwttoken.ResponseJWTUser{}, fmt.Errorf("%s field create JWT token: %w", op, ErrServic)
	}

	log.Info("success authentication user")

	return token, nil
}

func (c *UserServer) ServerDeleteUser(us domain.DeleteUser) error {
	const op = "user.ServerDeleteUser"

	log := c.log.WithField("op", op)

	log.Info("start delete user")

	if err := validator.New().Struct(us); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return fmt.Errorf("%s invalid validate: %w", op, err)
	}

	err := c.d.DeleteUser(us.Email, us.Password)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("not found")
			return fmt.Errorf("%s not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("field delete user")
		return fmt.Errorf("%s field delete user: %w", op, ErrDatabase)
	}

	log.Info("success delete user")
	return nil
}
