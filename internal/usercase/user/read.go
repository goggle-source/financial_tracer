package user

import (
	"fmt"

	"github.com/financial_tracer/internal/domain/entities"
)

type GetUserRepository interface {
	GetUser(email string, password string) (entities.User, error)
}

type GetUserServer struct {
	GetUserRepository
}

func (g *GetUserServer) ServerGetUser(email string, password string) (entities.User, error) {
	const op = "server.ServerGetUser"

	resultUser, err := g.GetUser(email, password)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return resultUser, nil
}
