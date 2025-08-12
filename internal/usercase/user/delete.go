package user

import (
	"fmt"

	"github.com/financial_tracer/internal/domain/entities"
)

type DeleteUserRepository interface {
	DeleteUser(email string, password string) error
}

type DeleteUserServer struct {
	DeleteUserRepository
}

func (d *DeleteUserServer) ServerDeleteUser(email string, password string) error {
	const op = "server.ServerDeleteUser"

	if entities.ValidEmail(email) {
		return entities.ErrorValidEmail
	}

	err := d.DeleteUser(email, password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
