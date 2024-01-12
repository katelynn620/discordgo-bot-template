package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string
}
