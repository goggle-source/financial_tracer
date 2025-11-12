package userHandlers

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	jwttoken "github.com/financial_tracer/internal/lib/jwtToken"
	"github.com/stretchr/testify/mock"
)

const secretKey = "secret"

type userServiceMock struct {
	mock.Mock
}

func (m *userServiceMock) RegistrationUser(ctx context.Context, us domain.RegisterUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(ctx, us)
	response, err := jwttoken.PostJWT(secretKey, args.Get(0).(uint), us.Name)
	if err != nil {
		return jwttoken.ResponseJWTUser{}, args.Error(1)
	}
	return response, args.Error(1)
}
func (m *userServiceMock) AuthenticationUser(ctx context.Context, us domain.AuthenticationUser) (jwttoken.ResponseJWTUser, error) {
	args := m.Called(ctx, us)
	response, err := jwttoken.PostJWT(secretKey, args.Get(0).(uint), args.Get(1).(string))
	if err != nil {
		return jwttoken.ResponseJWTUser{}, args.Error(2)
	}
	return response, args.Error(2)
}

func (m *userServiceMock) DeleteUser(ctx context.Context, us domain.DeleteUser) error {
	args := m.Called(ctx, us)
	return args.Error(0)
}
