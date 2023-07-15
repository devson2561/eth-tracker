package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Hash  string
	Value string
	TxRaw string
	From  string
	To    string
}
