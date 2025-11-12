package user

import (
	"context"
	"errors"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerRegistrationUser(t *testing.T) {

	type test struct {
		name         string
		user         domain.RegisterUser
		userID       uint
		mokuErr      error
		userErr      error
		shouldCallDB bool
	}

	arrTests := []test{
		{
			name: "success",
			user: domain.RegisterUser{
				Name:     "jonnsina",
				Email:    "jonn12@gmail.com",
				Password: "fgpDIJGP:OGhiHG",
			},
			userID:       1,
			mokuErr:      nil,
			userErr:      nil,
			shouldCallDB: true,
		},
		{
			name: "error database",
			user: domain.RegisterUser{
				Name:     "Sashasss",
				Email:    "tanks1235@gmail.com",
				Password: "adsafsgfGDgSD",
			},
			userID:       0,
			mokuErr:      errors.New("error database"),
			userErr:      ErrDatabase,
			shouldCallDB: true,
		},
		{
			name: "error valdiate",
			user: domain.RegisterUser{
				Name:     "",
				Email:    "afdaSFASF",
				Password: "asdasdadas",
			},
			userID:       0,
			mokuErr:      nil,
			userErr:      validator.ValidationErrors{},
			shouldCallDB: false,
		},
		{
			name: "error duplicated",
			user: domain.RegisterUser{
				Name:     "Sashasss",
				Email:    "tanks1235@gmail.com",
				Password: "adsafsgfGDgSD",
			},
			userID:       0,
			mokuErr:      postgresql.ErrorDuplicated,
			userErr:      ErrDuplicated,
			shouldCallDB: true,
		},
	}

	for _, test := range arrTests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("RegistrationUser", mock.Anything, mock.AnythingOfType("domain.User")).Return(test.userID, test.user.Name, test.mokuErr)

			server := CreateUserServer(repoMock, repoMock, repoMock, "secret", log)
			tokens, err := server.RegistrationUser(context.Background(), test.user)

			if test.mokuErr != nil || test.userErr != nil {
				assert.Error(t, err)
				if test.name == "error valdiate" {
					var ValidationErrors validator.ValidationErrors

					if !errors.As(err, &ValidationErrors) {
						t.Error("err != validator.ValidatiobErrors")
					}
				} else {
					if !errors.Is(err, test.userErr) {
						t.Error("err != test.userErr")
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "RegistrationUser", mock.Anything, mock.AnythingOfType("domain.User"))
			}
		})
	}
}

func TestServerAuthenticationUser(t *testing.T) {
	type test struct {
		name         string
		inputUser    domain.AuthenticationUser
		userID       uint
		nameUser     string
		mokuErr      error
		userErr      error
		shouldCallDB bool
	}

	tests := []test{
		{
			name: "success",
			inputUser: domain.AuthenticationUser{
				Email:    "jonn1342@gmail.com",
				Password: "admin12241532",
			},
			userID:       3,
			nameUser:     "jonn",
			mokuErr:      nil,
			userErr:      nil,
			shouldCallDB: true,
		},
		{
			name: "error database",
			inputUser: domain.AuthenticationUser{
				Email:    "sasha@gmail.com",
				Password: "admin15412",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      errors.New("error database"),
			userErr:      ErrDatabase,
			shouldCallDB: true,
		},
		{
			name: "error valdiate",
			inputUser: domain.AuthenticationUser{
				Email:    "aGfSGpGD",
				Password: "no password",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      nil,
			userErr:      validator.ValidationErrors{},
			shouldCallDB: false,
		},
		{
			name: "not found",
			inputUser: domain.AuthenticationUser{
				Email:    "jonn@gmail.com",
				Password: "12345678",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      postgresql.ErrorNotFound,
			userErr:      ErrNoFound,
			shouldCallDB: true,
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("AuthenticationUser", mock.Anything, ts.inputUser.Email, ts.inputUser.Password).
				Return(ts.userID, ts.nameUser, ts.mokuErr)

			server := CreateUserServer(repoMock, repoMock, repoMock, "secret", log)
			tokens, err := server.AuthenticationUser(context.Background(), ts.inputUser)

			if ts.mokuErr != nil || ts.userErr != nil {
				assert.Error(t, err)
				if ts.name == "error valdiate" {
					var ValidationErrors validator.ValidationErrors

					if !errors.As(err, &ValidationErrors) {
						t.Error("err != validator.ValidatiobErrors")
					}
				} else {
					if !errors.Is(err, ts.userErr) {
						t.Error("err != test.userErr")
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}

			if ts.shouldCallDB {
				repoMock.AssertCalled(t, "AuthenticationUser", mock.Anything, ts.inputUser.Email, ts.inputUser.Password)
			}
		})
	}
}

func TestServerDeleteUser(t *testing.T) {
	type test struct {
		name         string
		user         domain.DeleteUser
		mockErr      error
		userErr      error
		shouldCallDB bool
	}

	arrTests := []test{
		{
			name: "success",
			user: domain.DeleteUser{
				Email:    "jonn@gmail.com",
				Password: "admin",
			},
			mockErr:      nil,
			userErr:      nil,
			shouldCallDB: true,
		},
		{
			name: "error database",
			user: domain.DeleteUser{
				Email:    "sasha@gmail.com",
				Password: "admin15412",
			},
			mockErr:      errors.New("error database"),
			userErr:      ErrDatabase,
			shouldCallDB: true,
		},
		{
			name: "error not found",
			user: domain.DeleteUser{
				Email:    "jonn11@gmail.com",
				Password: "adminov",
			},
			mockErr:      postgresql.ErrorNotFound,
			userErr:      ErrNoFound,
			shouldCallDB: true,
		},
		{
			name: "error valdiate",
			user: domain.DeleteUser{
				Email:    "jonn@gmail.com",
				Password: "",
			},
			mockErr:      nil,
			userErr:      validator.ValidationErrors{},
			shouldCallDB: false,
		},
	}

	for _, ts := range arrTests {
		t.Run(ts.name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("DeleteUser", mock.Anything, ts.user.Email, ts.user.Password).Return(ts.mockErr)

			server := CreateUserServer(repoMock, repoMock, repoMock, "secret", log)
			err := server.DeleteUser(context.Background(), ts.user)

			if ts.mockErr != nil || ts.userErr != nil {
				assert.Error(t, err)
				if ts.name == "error valdiate" {
					var ValidationErrors validator.ValidationErrors

					if !errors.As(err, &ValidationErrors) {
						t.Error("err != validator.ValidatiobErrors")
					}
				} else {
					if !errors.Is(err, ts.userErr) {
						t.Error("err != test.userErr")
					}
				}
			} else {
				assert.NoError(t, err)
			}

			if ts.shouldCallDB {
				repoMock.AssertCalled(t, "DeleteUser", mock.Anything, ts.user.Email, ts.user.Password)
			}
		})
	}
}
