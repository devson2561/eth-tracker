package repository

import (
	"github.com/devson2561/eth-tracker/models"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *models.Address) error
	Find(address string) (*models.Address, error)
	FindAll() ([]models.Address, error)
}

type addressRepository struct {
	repo Repository
}

func NewAddressRepository(repo Repository) AddressRepository {
	return &addressRepository{
		repo: repo,
	}
}

func (r *addressRepository) Create(address *models.Address) error {
	return r.repo.Create(address)
}

func (r *addressRepository) FindAll() ([]models.Address, error) {
	var addresses []models.Address
	err := r.repo.FindAll(&addresses)
	return addresses, err
}

func (r *addressRepository) Find(address string) (*models.Address, error) {
	var a models.Address
	if err := r.repo.Find(&a, "address = ?", address); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}
