package repository

import (
	"boiler-plate/internal/settings/domain"
	baseModel "boiler-plate/pkg/db"
	"boiler-plate/pkg/errs"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repo struct {
	db   *gorm.DB
	base *baseModel.SQLClientRepository
}

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) Repository {
	return &Repo{db: db, base: base}
}

func (r Repo) FindSettings(ctx context.Context) (*domain.MainTable, error) {
	var (
		models *domain.MainTable
	)
	if err := r.db.WithContext(ctx).
		Model(&domain.MainTable{}).
		First(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models, nil
		}
		return nil, errs.Wrap(err)
	}
	return models, nil
}
