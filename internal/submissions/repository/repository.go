package repository

import (
	"boiler-plate/internal/submissions/domain"
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

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) SubmissionsRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.Submissions) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.Submissions{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Delete(ctx context.Context, tx *gorm.DB, key int) error {
	query := tx.WithContext(ctx)

	if err := query.
		Delete(&domain.Submissions{ID: key}).Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Submissions) error {
	query := tx.WithContext(ctx)
	if err := query.
		Model(&domain.Submissions{ID: id}).
		Updates(model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Find(ctx context.Context, tx *gorm.DB, limit, page int) (
	*[]domain.Submissions, *baseModel.Paginate, error,
) {
	var (
		models *[]domain.Submissions
	)
	tx = tx.WithContext(ctx).
		Model(&domain.Submissions{})
	pagination := baseModel.NewPaginate(limit, page)
	if err := tx.
		Scopes(pagination.PaginatedResult(&models, tx)).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pagination, nil
		}
		return nil, nil, errs.Wrap(err)
	}
	return models, pagination, nil
}

func (r Repo) FindByUser(ctx context.Context, tx *gorm.DB, limit, page, userid int) (
	*[]domain.Submissions, *baseModel.Paginate, error,
) {
	var (
		models *[]domain.Submissions
	)
	pagination := baseModel.NewPaginate(limit, page)
	tx = tx.WithContext(ctx).
		Model(&domain.Submissions{}).
		Where("user_id = ?", userid)
	if err := tx.
		Scopes(pagination.PaginatedResult(&models, tx)).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pagination, nil
		}
		return nil, nil, errs.Wrap(err)
	}
	return models, pagination, nil
}

func (r Repo) Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Submissions, error) {
	var (
		models *domain.Submissions
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Submissions{}).
		Preload("User").
		First(&models, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}

func (r Repo) DetailByUser(ctx context.Context, tx *gorm.DB, id int) (*domain.Submissions, error) {
	var (
		models *domain.Submissions
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Submissions{}).
		Order("created_at DESC").
		First(&models, "user_id = ?", id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}
