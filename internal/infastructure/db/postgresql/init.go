package postgresql

import (
	"fmt"

	"github.com/financial_tracer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name         string `gorm:"size:50;not null"`
	Email        string `gorm:"not null;enique"`
	PasswordHash []byte `gorm:"not null"`
}

type Db struct {
	DB *gorm.DB
}

func Init(cfg *config.Config) (*Db, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		cfg.App.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.PortDb, cfg.DB.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error conn database: %w", err)
	}

	err = db.AutoMigrate(
		&Users{},
	)
	if err != nil {
		return nil, fmt.Errorf("error migrate database: %w", err)
	}

	return &Db{
		DB: db,
	}, nil
}
