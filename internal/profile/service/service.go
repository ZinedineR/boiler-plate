package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/profile/domain"
	"boiler-plate/internal/profile/repository"
	"boiler-plate/pkg/exception"
	"context"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strconv"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.ProfileRepository, db *gorm.DB, validate *validator.Validate,
) Service {
	return &service{config: config, ProfileRepo: repo, validate: validate, DB: db}
}

type service struct {
	DB          *gorm.DB
	config      *appconf.Config
	ProfileRepo repository.ProfileRepository
	validate    *validator.Validate
}

func (s service) Create(
	ctx context.Context, req *domain.Profile,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	err := s.ProfileRepo.Create(ctx, tx, req)
	if err != nil {
		exception.Internal("error inserting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Update(
	ctx context.Context, id string, profile *domain.Profile,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return exception.PermissionDenied("Input of id must be integer")
	}
	err = s.ProfileRepo.Update(ctx, tx, idInt, profile)
	if err != nil {
		exception.Internal("error inserting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Delete(ctx context.Context, id string) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return exception.PermissionDenied("Input of id must be integer")
	}
	err = s.ProfileRepo.Delete(ctx, tx, idInt)
	if err != nil {
		exception.Internal("error deleting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Find(ctx context.Context) (*[]domain.Profile, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	result, err := s.ProfileRepo.Find(ctx, tx)
	if err != nil {
		return nil, exception.Internal("error geting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}

func (s service) Detail(ctx context.Context, id string) (*domain.Profile, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, exception.PermissionDenied("Input of id must be integer")
	}
	result, err := s.ProfileRepo.Detail(ctx, tx, idInt)
	if err != nil {
		return nil, exception.Internal("error deleting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}

func (s service) Auth(ctx context.Context, profile, password string) (*domain.Profile, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	result, err := s.ProfileRepo.Auth(ctx, tx, profile, password)
	if err != nil {
		return nil, exception.Internal("error deleting profile", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
