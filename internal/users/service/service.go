package service

import (
	"boiler-plate/app/appconf"
	subRepo "boiler-plate/internal/submissions/repository"
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/repository"
	"boiler-plate/pkg/exception"
	"context"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strconv"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.UsersRepository, submissionrepo subRepo.SubmissionsRepository, db *gorm.DB,
	validate *validator.Validate,
) Service {
	return &service{config: config, UsersRepo: repo, SubmissionsRepo: submissionrepo, validate: validate, DB: db}
}

type service struct {
	DB              *gorm.DB
	config          *appconf.Config
	UsersRepo       repository.UsersRepository
	SubmissionsRepo subRepo.SubmissionsRepository
	validate        *validator.Validate
}

func (s service) Create(
	ctx context.Context, req *domain.Users,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if err := s.validate.Struct(req); err != nil {
		return exception.InvalidArgument(err)
	}
	err := s.UsersRepo.Create(ctx, tx, req)
	if err != nil {
		return exception.Internal("error inserting users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Update(
	ctx context.Context, id string, users *domain.Users,
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
	err = s.UsersRepo.Update(ctx, tx, idInt, users)
	if err != nil {
		return exception.Internal("error updating users", err)
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
	if len(*result) > 0 {
		for i, response := range *result {
			latestSubmissions, err := s.SubmissionsRepo.DetailByUser(ctx, tx, response.ID)
			if err != nil {
				return nil, exception.Internal("error getting detail submissions", err)
			}
			if latestSubmissions != nil {
				(*result)[i].DeclareRiskProfile(latestSubmissions.RiskScore, latestSubmissions.RiskCategory, latestSubmissions.RiskDefinition)
			}
		}
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

func (s service) Detail(ctx context.Context, id string) (*domain.UserResponse, *exception.Exception) {
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
	latestSubmissions, err := s.SubmissionsRepo.DetailByUser(ctx, tx, result.ID)
	if err != nil {
		return nil, exception.Internal("error getting detail submissions", err)
	}
	if latestSubmissions != nil {
		result.DeclareRiskProfile(latestSubmissions.RiskScore, latestSubmissions.RiskCategory, latestSubmissions.RiskDefinition)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}

func (s service) Auth(ctx context.Context, email, password string) (*domain.Users, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	result, err := s.UsersRepo.Auth(ctx, tx, email, password)
	if err != nil {
		return nil, exception.Internal("error finding users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
