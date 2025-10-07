package user

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) RegistrationUser(user domain.User) (uint, string, error) {
	args := d.Called(user)
	return args.Get(0).(uint), args.String(1), args.Error(2)
}

func (d *DbMock) DeleteUser(email string, password string) error {
	args := d.Called(email, password)
	return args.Error(0)
}

func (d *DbMock) AuthenticationUser(email string, password string) (uint, string, error) {
	args := d.Called(email, password)
	return args.Get(0).(uint), args.String(1), args.Error(2)
}
