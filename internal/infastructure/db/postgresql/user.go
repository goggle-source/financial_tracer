package postgresql

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (d *Db) RegistrationUser(user *domain.User) (int, string, error) {
	const op = "Postgresql.Login"

	userDb := Users{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash[:],
	}

	result := d.DB.Create(&userDb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, "", fmt.Errorf("%s: %w", op, domain.ErrorDuplicated)
		}
		return 0, "", fmt.Errorf("%s: %w", op, result.Error)
	}
	return int(userDb.ID), user.Name, nil
}

func (d *Db) AuthenticationUser(email string, password string) (int, string, error) {
	const op = "postgresql.AuthenticationUser"

	var user Users

	result := d.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, "", fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return 0, "", fmt.Errorf("%s: %w", op, result.Error)
	}

	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return int(user.ID), user.Name, nil
}

func (d *Db) DeleteUser(email string, passwordHash string) error {
	const op = "postgresql.DeleteUser"
	var user domain.User
	result := d.DB.Where("email = ?", email, passwordHash).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passwordHash))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result = d.DB.Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}
