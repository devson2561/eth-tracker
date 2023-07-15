package repository

import (
	"github.com/devson2561/eth-tracker/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	Find(hash string) (*models.Transaction, error)
	FindAll() ([]models.Transaction, error)
}

type transactionRepository struct {
	repo Repository
}

func NewTransactionRepository(repo Repository) TransactionRepository {
	return &transactionRepository{
		repo: repo,
	}
}

func (r *transactionRepository) Create(address *models.Transaction) error {
	return r.repo.Create(address)
}

func (r *transactionRepository) Find(hash string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.repo.Find(&transaction, "hash = ?", hash); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.repo.FindAll(&transactions)
	return transactions, err
}
