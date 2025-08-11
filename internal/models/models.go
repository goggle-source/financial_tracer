package models

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:50;not null"`
	Email     string    `gorm:"size:100;not null;enique"`
	Password  string    `gorm:"size:200;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:create_user"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:update_user"`
}
