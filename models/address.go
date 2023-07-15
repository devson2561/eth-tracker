package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Address string `gorm:"unique"`
}
