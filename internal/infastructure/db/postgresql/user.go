package postgresql

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/models"
	"gorm.io/gorm"
)

func (d *Db) RegistrationUser(user *domain.User) error {
	const op = "Postgresql.Login"

	userDb := Users{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash[:],
	}

	result := d.DB.Create(&userDb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%s: %w", op, domain.ErrorDuplicated)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}

func (d *Db) AuthenticationUser(email string, password []byte) (*models.UserResponse, error) {
	const op = "postgresql.GetUser"

	var user Users

	result := d.DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	userDomain := models.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}

	return &userDomain, nil
}

func (d *Db) DeleteUser(email string, passwordHash []byte) error {
	const op = "postgresql.DeleteUser"
	var user domain.User
	result := d.DB.Where("email = ? AND password = ?", email, passwordHash).Delete(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrorNotFound)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}
	return nil
}
