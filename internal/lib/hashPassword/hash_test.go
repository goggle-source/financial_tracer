package hashPassword

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	arr := []string{"jonnn", "sasha", "godsIGds"}
	for _, password := range arr {
		t.Run("hash", func(t *testing.T) {
			hash, err := Hash(password)
			if err != nil {
				t.Error("error create hash")
			}
			err = bcrypt.CompareHashAndPassword(hash, []byte(password))
			if err != nil {
				t.Error("error: hash is not equal is password")
			}
		})
	}
}
