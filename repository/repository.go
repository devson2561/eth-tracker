package repository

import "gorm.io/gorm"

type Repository interface {
	Create(value interface{}) error
	Find(out interface{}, where ...interface{}) error
	FindAll(out interface{}, where ...interface{}) error
	Delete(value interface{}) error
	OrderFirst(out interface{}, order string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{
		db: db,
	}
}

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) Create(value interface{}) error {
	return r.db.Create(value).Error
}

func (r *gormRepository) Find(out interface{}, where ...interface{}) error {
	return r.db.First(out, where...).Error
}

func (r *gormRepository) FindAll(out interface{}, where ...interface{}) error {
	return r.db.Find(out, where...).Error
}

func (r *gormRepository) Delete(value interface{}) error {
	return r.db.Delete(value).Error
}

func (r *gormRepository) OrderFirst(out interface{}, order string) error {
	return r.db.Order(order).First(out).Error
}
