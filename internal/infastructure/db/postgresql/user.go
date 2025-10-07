package postgresql

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (d *Db) RegistrationUser(user domain.User) (uint, string, error) {

	userDb := User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash[:],
	}

	result := d.DB.Create(&userDb)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, "", ErrorDuplicated
		}
		return 0, "", result.Error
	}
	return userDb.ID, user.Name, nil
}

func (d *Db) AuthenticationUser(email string, password string) (uint, string, error) {

	var user User

	result := d.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, "", ErrorNotFound
		}
		return 0, "", result.Error
	}

	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return 0, "", err
	}

	return user.ID, user.Name, nil
}

func (d *Db) DeleteUser(email string, passwordHash string) error {
	var user domain.User
	result := d.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrorNotFound
		}
		return result.Error
	}
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passwordHash))
	if err != nil {
		return err
	}

	result = d.DB.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}
