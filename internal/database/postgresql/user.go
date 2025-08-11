package postgresql

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/database/storage"
	"github.com/financial_tracer/internal/models"
	"gorm.io/gorm"
)

func (d *Db) CreateUser(name string, email string, password string) (uint, error) {
	const op = "Postgresql.Login"
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	result := d.DB.Create(&user)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}
	return user.Id, nil
}

func (d *Db) GetUser(name string, email string, password string) (uint, error) {
	const op = "postgresql.GetUser"

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	result := d.DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrorNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user.Id, nil
}

func (d *Db) DeleteUser(email string, password string) error {
	const op = "postgresql.DeleteUser"
	var user models.User
	result := d.DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, storage.ErrorNotFound)
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	err := d.DB.Select("Categories", "Transactions").Delete(&user).Error
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
