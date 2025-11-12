package user

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) RegistrationUser(ctx context.Context, user domain.User) (uint, string, error) {
	args := d.Called(ctx, user)
	return args.Get(0).(uint), args.String(1), args.Error(2)
}

func (d *DbMock) DeleteUser(ctx context.Context, email string, password string) error {
	args := d.Called(ctx, email, password)
	return args.Error(0)
}

func (d *DbMock) AuthenticationUser(ctx context.Context, email string, password string) (uint, string, error) {
	args := d.Called(ctx, email, password)
	return args.Get(0).(uint), args.String(1), args.Error(2)
}
