package postgresql

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain/entities"
	storage "github.com/financial_tracer/internal/domain/errors"
	"gorm.io/gorm"
)

func (d *Db) CreateUser(user *entities.User) error {
	const op = "Postgresql.Login"

	result := d.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%s: %w", op, storage.ErrorDuplicated)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (d *Db) GetUser(email string, password string) (*entities.User, error) {
	const op = "postgresql.GetUser"

	var user entities.User

	result := d.DB.Where("email = ? AND password = ?", user.Email, user.PasswordHash).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &entities.User{}, fmt.Errorf("%s: %w", op, storage.ErrorNotFound)
		}
		return &entities.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &user, nil
}

func (d *Db) DeleteUser(email string, password string) error {
	const op = "postgresql.DeleteUser"
	var user entities.User
	result := d.DB.Where("email = ? AND password = ?", email, password).Delete(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, storage.ErrorNotFound)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}
