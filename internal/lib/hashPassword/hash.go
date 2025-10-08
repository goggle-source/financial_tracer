package hashPassword

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), fmt.Errorf("error hash password")
	}
	return bytes, nil
}
