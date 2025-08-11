package postgresql

import (
	"fmt"

	"github.com/financial_tracer/internal/config"
	"github.com/financial_tracer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		&models.User{},
	)
	if err != nil {
		return nil, fmt.Errorf("error migrate database: %w", err)
	}

	return &Db{
		DB: db,
	}, nil
}
