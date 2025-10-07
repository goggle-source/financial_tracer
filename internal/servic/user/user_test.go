package user

import (
	"errors"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
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
		msgErr       string
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
			msgErr:       "",
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
			msgErr:       "field registration user",
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
			msgErr:       "invalid validate",
			shouldCallDB: false,
		},
	}

	for _, test := range arrTests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("RegistrationUser", mock.AnythingOfType("domain.User")).Return(test.userID, test.user.Name, test.mokuErr)

			server := CreateUserServer(repoMock, "secret", log)
			tokens, err := server.ServerRegistrationUser(test.user)

			if test.msgErr != "" || test.mokuErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.msgErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "RegistrationUser", mock.AnythingOfType("domain.User"))
			}
		})
	}
}

func TestServerAuthenticationUser(t *testing.T) {
	type test struct {
		Name         string
		inputUser    domain.AuthenticationUser
		userID       uint
		nameUser     string
		mokuErr      error
		msgErr       string
		shouldCallDB bool
	}

	tests := []test{
		{
			Name: "success",
			inputUser: domain.AuthenticationUser{
				Email:    "jonn1342@gmail.com",
				Password: "admin12241532",
			},
			userID:       3,
			nameUser:     "jonn",
			mokuErr:      nil,
			msgErr:       "",
			shouldCallDB: true,
		},
		{
			Name: "error database",
			inputUser: domain.AuthenticationUser{
				Email:    "sasha@gmail.com",
				Password: "admin15412",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      errors.New("error database"),
			msgErr:       "field authentication user",
			shouldCallDB: true,
		},
		{
			Name: "error validate",
			inputUser: domain.AuthenticationUser{
				Email:    "aGfSGpGD",
				Password: "no password",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      nil,
			msgErr:       "invalid validate",
			shouldCallDB: false,
		},
		{
			Name: "not found",
			inputUser: domain.AuthenticationUser{
				Email:    "jonn@gmail.com",
				Password: "12345678",
			},
			userID:       0,
			nameUser:     "",
			mokuErr:      postgresql.ErrorNotFound,
			msgErr:       "field authentication user",
			shouldCallDB: true,
		},
	}

	for _, ts := range tests {
		t.Run(ts.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("AuthenticationUser", ts.inputUser.Email, ts.inputUser.Password).
				Return(ts.userID, ts.nameUser, ts.mokuErr)

			server := CreateUserServer(repoMock, "secret", log)
			tokens, err := server.ServerAuthenticationUser(ts.inputUser)

			if ts.mokuErr != nil || ts.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), ts.msgErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}

			if ts.shouldCallDB {
				repoMock.AssertCalled(t, "AuthenticationUser", ts.inputUser.Email, ts.inputUser.Password)
			}
		})
	}
}

func TestServerDeleteUser(t *testing.T) {
	type test struct {
		Name         string
		user         domain.DeleteUser
		mockErr      error
		msgErr       string
		shouldCallDB bool
	}

	arrTests := []test{
		{
			Name: "success",
			user: domain.DeleteUser{
				Email:    "jonn@gmail.com",
				Password: "admin",
			},
			mockErr:      nil,
			msgErr:       "",
			shouldCallDB: true,
		},
		{
			Name: "error database",
			user: domain.DeleteUser{
				Email:    "sasha@gmail.com",
				Password: "admin15412",
			},
			mockErr:      errors.New("error database"),
			msgErr:       "field delete user",
			shouldCallDB: true,
		},
		{
			Name: "error not found",
			user: domain.DeleteUser{
				Email:    "jonn11@gmail.com",
				Password: "adminov",
			},
			mockErr:      postgresql.ErrorNotFound,
			msgErr:       "field delete user",
			shouldCallDB: true,
		},
		{
			Name: "validate",
			user: domain.DeleteUser{
				Email:    "jonn@gmail.com",
				Password: "",
			},
			mockErr:      nil,
			msgErr:       "invalid validate",
			shouldCallDB: false,
		},
	}

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			log := logrus.New()

			repoMock.On("DeleteUser", test.user.Email, test.user.Password).Return(test.mockErr)

			server := CreateUserServer(repoMock, "secret", log)
			err := server.ServerDeleteUser(test.user)

			if test.mockErr != nil || test.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.msgErr)
			} else {
				assert.NoError(t, err)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "DeleteUser", test.user.Email, test.user.Password)
			}
		})
	}
}
