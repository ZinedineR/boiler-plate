package repository

import (
	"boiler-plate/internal/transaction/domain"
	baseModel "boiler-plate/pkg/db"
	"boiler-plate/pkg/errs"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repo struct {
	db   *gorm.DB
	base *baseModel.SQLClientRepository
}

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) TransactionRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.Transaction) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.Transaction{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Transaction) error {
	query := tx.WithContext(ctx)
	if err := query.
		Model(&domain.Transaction{ID: id}).
		Updates(model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Find(ctx context.Context, tx *gorm.DB, limit, page, userid int) (
	*[]domain.Transaction, *baseModel.Paginate, error,
) {
	var (
		models *[]domain.Transaction
	)
	tx = tx.WithContext(ctx).Where("user_id = ?", userid).
		Model(&domain.Transaction{})
	pagination := baseModel.NewPaginate(limit, page)
	if err := tx.
		Scopes(pagination.PaginatedResult(&domain.Transaction{}, tx)).
		Find(&models).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pagination, nil
		}
		return nil, nil, errs.Wrap(err)
	}
	return models, pagination, nil
}

func (r Repo) Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Transaction, error) {
	var (
		models *domain.Transaction
	)

	if err := tx.WithContext(ctx).
		Model(&domain.Transaction{}).
		First(&models, id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return models, nil
}
