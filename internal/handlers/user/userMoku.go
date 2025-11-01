package userHandlers

import (
	"github.com/financial_tracer/internal/domain"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/stretchr/testify/mock"
)

const secretKey = "secret"

type userServiceMock struct {
	mock.Mock
}

func (m *userServiceMock) RegistrationUser(us domain.RegisterUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(us)
	response, err := jwttoken.PostJWT(secretKey, args.Get(0).(uint), us.Name)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, args.Error(1)
	}
	return response, args.Error(1)
}
func (m *userServiceMock) AuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(us)
	response, err := jwttoken.PostJWT(secretKey, args.Get(0).(uint), args.Get(1).(string))
	if err != nil {
		return jwttoken.ResponseJWTUser{}, args.Error(2)
	}
	return response, args.Error(2)
}

func (m *userServiceMock) DeleteUser(us domain.DeleteUser) error {
	args := m.Called(us)
	return args.Error(0)
}
