package userHandlers

import (
	"github.com/financial_tracer/internal/domain"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	mock.Mock
}

func (m *userServiceMock) ServerRegistrationUser(us domain.RegisterUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(us)
	return args.Get(0).(jwttoken.ResponseJWTUser), args.Error(2)
}
func (m *userServiceMock) ServerAuthenticationUser(us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(us)
	return args.Get(0).(jwttoken.ResponseJWTUser), args.Error(2)
}
func (m *userServiceMock) ServerDeleteUser(us domain.DeleteUser) error {
	args := m.Called(us)
	return args.Error(0)
}
