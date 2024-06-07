package repository

import (
	"boiler-plate/internal/profile/domain"
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

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) ProfileRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.Profile) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.Profile{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Delete(ctx context.Context, tx *gorm.DB, key int) error {
	query := tx.WithContext(ctx)

	if err := query.
		Delete(&domain.Profile{ID: key}).Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Profile) error {
	query := tx.WithContext(ctx)
	if err := query.
		Model(&domain.Profile{ID: id}).
		Updates(model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Find(ctx context.Context, tx *gorm.DB) (*[]domain.Profile, error) {
	var (
		models *[]domain.Profile
	)
	if err := tx.WithContext(ctx).
		Model(&domain.Profile{}).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models, nil
		}
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (r Repo) Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Profile, error) {
	var (
		models *domain.Profile
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Profile{}).
		First(&models, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}

func (r Repo) Auth(ctx context.Context, tx *gorm.DB, profile, password string) (*domain.Profile, error) {
	var (
		models *domain.Profile
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("profile = ?", profile).Where("password = ?", password).
		First(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}
