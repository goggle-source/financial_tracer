package user

import (
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
)

type DatabaseUserRepository interface {
	RegistrationUser(user *domain.User) error
	DeleteUser(email string, passwordHash []byte) error
	AuthenticationUser(email string, password []byte) (*models.UserResponse, error)
}

type CreateUserServer struct {
	d DatabaseUserRepository
}

func CreateServer(d DatabaseUserRepository) *CreateUserServer {
	return &CreateUserServer{
		d: d,
	}
}

const lenRequestId = 100000

func (c *CreateUserServer) ServerRegistrationUser(userRequest models.UserRequest) (*models.UserResponse, error) {
	const op = "server.ServerCreateUser"

	passwordHash, err := Hash(userRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}
	requestId := NewRandomString(lenRequestId)

	user := domain.User{
		Name:         userRequest.Name,
		Email:        userRequest.Email,
		RequestId:    requestId,
		PasswordHash: passwordHash,
	}

	err = c.d.RegistrationUser(&user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (c *CreateUserServer) ServerAuthenticationUser(email string, password string) (*models.UserResponse, error) {
	const op = "server.ServerGetUser"

	passwordHash, err := Hash(password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	resultUser, err := c.d.AuthenticationUser(email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.UserResponse{
		Name:  resultUser.Name,
		Email: resultUser.Email,
	}, nil
}

func (c *CreateUserServer) ServerDeleteUser(email string, password string) error {
	const op = "server.ServerDeleteUser"

	passwordHash, err := Hash(password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, domain.ErrorHashPassword)
	}

	err = c.d.DeleteUser(email, passwordHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
