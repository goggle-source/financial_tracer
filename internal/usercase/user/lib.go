package user

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), fmt.Errorf("error hash password")
	}
	return bytes, nil
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
