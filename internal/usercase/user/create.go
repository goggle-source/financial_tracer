package user

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/financial_tracer/internal/domain/entities"
)

type CreateUserRepository interface {
	CreateUser(*entities.User) error
}

type CreateUserServer struct {
	CreateUserRepository
}

const lenRequestId = 10000

func (c *CreateUserServer) ServerCreateUser(name string, email string, password string) (*entities.User, error) {
	const op = "server.ServerCreateUser"

	passwordHash := sha256.Sum256([]byte(password))
	requestId := NewRandomString(lenRequestId)

	user, err := entities.NewUser(requestId, name, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = c.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	arr := []rune("qwertyuiopasdfghjklzxcvbnm" + "QWERTYUIOPASDFGHJKLZXCVBNM" + "1234567890" + "*#$")

	result := make([]rune, size)

	for i := range result {
		result[i] = arr[rnd.Intn(len(arr))]
	}

	return string(result)
}
