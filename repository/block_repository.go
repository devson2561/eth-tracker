package repository

import (
	"github.com/devson2561/eth-tracker/models"
	"gorm.io/gorm"
)

type BlockRepository interface {
	Create(block *models.Block) error
	Find(blockNumber uint64) (*models.Block, error)
	FindAll() ([]models.Block, error)
	FindLatest() (*models.Block, error)
}

type blockRepository struct {
	repo Repository
}

func NewBlockRepository(repo Repository) BlockRepository {
	return &blockRepository{
		repo: repo,
	}
}

func (r *blockRepository) Create(block *models.Block) error {
	return r.repo.Create(block)
}

func (r *blockRepository) Find(blockNumber uint64) (*models.Block, error) {
	var block models.Block
	if err := r.repo.Find(&block, "block_number = ?", blockNumber); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

func (r *blockRepository) FindAll() ([]models.Block, error) {
	var blocks []models.Block
	err := r.repo.FindAll(&blocks)
	return blocks, err
}

func (r *blockRepository) FindLatest() (*models.Block, error) {
	var block models.Block
	err := r.repo.OrderFirst(&block, "block_number desc")
	if err != nil {
		return nil, err
	}
	return &block, nil
}
