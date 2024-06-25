package repository

import (
	"boiler-plate/internal/users/domain"
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

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) UsersRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.Users) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.Users{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Delete(ctx context.Context, tx *gorm.DB, key int) error {
	query := tx.WithContext(ctx)

	if err := query.
		Delete(&domain.Users{ID: key}).Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Users) error {
	query := tx.WithContext(ctx)
	if err := query.
		Model(&domain.Users{ID: id}).
		Updates(model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Find(ctx context.Context, tx *gorm.DB, limit, page int) (
	*[]domain.UserResponse, *baseModel.Paginate, error,
) {
	var (
		models *[]domain.UserResponse
	)
	tx = tx.WithContext(ctx).
		Select("id", "email", "password", "created_at", "updated_at").
		Model(&domain.Users{})
	pagination := baseModel.NewPaginate(limit, page)
	if err := tx.
		Scopes(pagination.PaginatedResult(&domain.Users{}, tx)).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pagination, nil
		}
		return nil, nil, errs.Wrap(err)
	}
	return models, pagination, nil
}

func (r Repo) Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.UserResponse, error) {
	var (
		models *domain.UserResponse
	)

	if err := tx.WithContext(ctx).
		Select("id", "email", "password", "created_at", "updated_at").
		Model(&domain.Users{}).
		First(&models, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}

func (r Repo) Auth(ctx context.Context, tx *gorm.DB, users, password string) (*domain.Users, error) {
	var (
		models *domain.Users
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Users{}).
		Where("email = ?", users).Where("password = ?", password).
		First(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}
