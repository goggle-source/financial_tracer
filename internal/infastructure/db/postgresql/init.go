package postgresql

import (
	"fmt"

	"github.com/financial_tracer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string        `gorm:"size:50;not null"`
	Email        string        `gorm:"not null;unique"`
	PasswordHash []byte        `gorm:"not null"`
	Categories   []Category    `gorm:"foreignKey:UserID"`
	Transactions []Transaction `gorm:"foreignKey:UserID"`
}

type Category struct {
	gorm.Model
	Name         string `gorm:"size:60;not null;unique"`
	UserID       uint
	Limit        int           `gorm:"not null"`
	Type         string        `gorm:"size:100"`
	Description  string        `gorm:"size:100"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID"`
}

type Transaction struct {
	gorm.Model
	Name        string `gorm:"not null;size:60"`
	UserID      uint
	CategoryID  uint
	Count       int    `gorm:"not null"`
	Description string `gorm:"size:100"`
}

type Db struct {
	DB *gorm.DB
}

func Init(cfg *config.Config) (*Db, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=%s",
		cfg.App.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.PortDb, cfg.DB.DbName, cfg.DB.Time)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error conn database: %w", err)
	}

	err = db.AutoMigrate(
		&User{},
		&Category{},
	)
	if err != nil {
		return nil, fmt.Errorf("error migrate database: %w", err)
	}

	return &Db{
		DB: db,
	}, nil
}
