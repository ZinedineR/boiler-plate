package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/repository"
	walletDomain "boiler-plate/internal/wallet/domain"
	walletRepo "boiler-plate/internal/wallet/repository"
	"boiler-plate/pkg/exception"
	"context"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strconv"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.UsersRepository, wallet walletRepo.WalletRepository, db *gorm.DB,
	validate *validator.Validate,
) Service {
	return &service{config: config, UsersRepo: repo, WalletRepo: wallet, validate: validate, DB: db}
}

type service struct {
	DB         *gorm.DB
	config     *appconf.Config
	UsersRepo  repository.UsersRepository
	WalletRepo walletRepo.WalletRepository
	validate   *validator.Validate
}

func (s service) Create(
	ctx context.Context, req *UserRequest,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if err := s.validate.Struct(req); err != nil {
		return exception.InvalidArgument(err)
	}
	walletBody := &walletDomain.Wallet{}
	if err := s.WalletRepo.Create(ctx, tx, walletBody); err != nil {
		return exception.Internal("failed to create wallet", err)
	}
	body := &domain.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		WalletId: walletBody.ID,
	}
	err := s.UsersRepo.Create(ctx, tx, body)
	if err != nil {
		return exception.Internal("error creating users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	req.Id = body.ID
	return nil
}

func (s service) Update(
	ctx context.Context, id string, users *UserRequest,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if err := s.validate.Struct(users); err != nil {
		return exception.InvalidArgument(err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return exception.PermissionDenied("Input of id must be integer")
	}
	body := &domain.Users{
		Name:     users.Name,
		Email:    users.Email,
		Password: users.Password,
	}
	err = s.UsersRepo.Update(ctx, tx, idInt, body)
	if err != nil {
		return exception.Internal("error updating users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	users.Id = body.ID
	return nil
}

func (s service) Delete(ctx context.Context, id string) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return exception.PermissionDenied("Input of id must be integer")
	}
	err = s.UsersRepo.Delete(ctx, tx, idInt)
	if err != nil {
		return exception.Internal("error deleting users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Find(ctx context.Context, limit string, page string) (*FindResponse, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	var limitInt, pageInt int

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, exception.PermissionDenied("Input of limit must be integer")
	}
	pageInt, err = strconv.Atoi(page)
	if err != nil {
		return nil, exception.PermissionDenied("Input of page must be integer")
	}
	result, pagination, err := s.UsersRepo.Find(ctx, tx, limitInt, pageInt)
	if err != nil {
		return nil, exception.Internal("error geting users", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	finalResponse := &FindResponse{
		Pagination: *pagination,
		Data:       *result,
	}
	return finalResponse, nil
}

func (s service) Detail(ctx context.Context, id string) (*domain.Users, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, exception.PermissionDenied("Input of id must be integer")
	}
	result, err := s.UsersRepo.Detail(ctx, tx, idInt)
	if err != nil {
		return nil, exception.Internal("error getting detail users", err)
	}
	if result == nil {
		return nil, exception.NotFound("detail not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
