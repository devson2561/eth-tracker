package models

import (
	"gorm.io/gorm"
)

type Block struct {
	gorm.Model
	BlockNumber      uint64 `gorm:"column:block_number"`
	BlockHash        string `gorm:"column:block_hash"`
	TransactionCount int    `gorm:"column:transaction_count"`
}
