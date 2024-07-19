package repository

import (
	"boiler-plate/internal/url/domain"
	baseModel "boiler-plate/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repo struct {
	db   *gorm.DB
	base *baseModel.SQLClientRepository
}

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) UsersRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.URL) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.URL{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Detail(ctx context.Context, tx *gorm.DB, longurl string) (*domain.URL, error) {
	var (
		models *domain.URL
	)

	if err := tx.WithContext(ctx).
		Model(&domain.URL{}).
		Where("shorturl = ?", longurl).
		First(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}
