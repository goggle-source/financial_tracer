package user

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/financial_tracer/internal/lib/hashPassword"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RegistrationuserRepository interface {
	RegistrationUser(user domain.User) (uint, string, error)
}

type DeleteUserRepository interface {
	DeleteUser(email string, password string) error
}

type AuthenticationUserRepository interface {
	AuthenticationUser(email string, password string) (uint, string, error)
}

type UserValid struct {
	Valid func(error) []validator.ValidationErrors
}

type UserServer struct {
	log       *logrus.Logger
	r         RegistrationuserRepository
	d         DeleteUserRepository
	a         AuthenticationUserRepository
	validate  validator.Validate
	SecretKey string
}

func CreateUserServer(r RegistrationuserRepository, d DeleteUserRepository, a AuthenticationUserRepository, sk string, log *logrus.Logger) *UserServer {
	return &UserServer{
		log:       log,
		d:         d,
		r:         r,
		a:         a,
		validate:  *validator.New(),
		SecretKey: sk,
	}
}

func (c *UserServer) RegistrationUser(us domain.RegisterUser) (jwttoken.ResponseJWTUser, error) {
	const op = "user.ServerRegistrationUser"

	log := c.log.WithField("op", op)

	log.Info("registration user")

	if err := c.validate.Struct(us); err != nil {
		log.WithFields(logrus.Fields{"err": err}).Error("invalid validate")

		return jwttoken.ResponseJWTUser{}, err
	}

	passwordHash, err := hashPassword.Hash(us.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("field hash password")

		return jwttoken.ResponseJWTUser{}, ErrServic
	}

	user := domain.User{
		Name:         us.Name,
		Email:        us.Email,
		PasswordHash: passwordHash,
	}

	id, name, err := c.r.RegistrationUser(user)
	if err != nil {
		if errors.Is(err, postgresql.ErrorDuplicated) {
			log.WithField("err", err).Error("field duplicated")

			return jwttoken.ResponseJWTUser{}, ErrDuplicated
		}
		log.WithField("err", err).Error("field registration user")
		return jwttoken.ResponseJWTUser{}, ErrDatabase
	}

	tokens, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		log.WithField("err", err).Error("field create JWT token")
		return jwttoken.ResponseJWTUser{}, ErrServic
	}

	log.Info("success registration user")

	return tokens, nil
}

func (c *UserServer) AuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	const op = "user.ServerAuthenticationUser"

	log := c.log.WithField("op", op)

	log.Info("start authentication user")

	if err := c.validate.Struct(us); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return jwttoken.ResponseJWTUser{}, err
	}

	id, name, err := c.a.AuthenticationUser(us.Email, us.Password)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("not found")

			return jwttoken.ResponseJWTUser{}, ErrNoFound
		}
		log.WithField("err", err).Error("field authentication user")
		return jwttoken.ResponseJWTUser{}, ErrDatabase
	}

	token, err := jwttoken.PostJWT(c.SecretKey, id, name)
	if err != nil {
		log.WithField("err", err).Error("field create JWT token")

		return jwttoken.ResponseJWTUser{}, ErrServic
	}

	log.Info("success authentication user")

	return token, nil
}

func (c *UserServer) DeleteUser(us domain.DeleteUser) error {
	const op = "user.ServerDeleteUser"

	log := c.log.WithField("op", op)

	log.Info("start delete user")

	if err := c.validate.Struct(us); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return err
	}

	err := c.d.DeleteUser(us.Email, us.Password)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("not found")
			return ErrNoFound
		}
		log.WithField("err", err).Error("field delete user")
		return ErrDatabase
	}

	log.Info("success delete user")
	return nil
}
