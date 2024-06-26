package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/internal/submissions/repository"
	"boiler-plate/pkg/exception"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strconv"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.SubmissionsRepository, db *gorm.DB, validate *validator.Validate,
) Service {
	return &service{config: config, SubmissionsRepo: repo, validate: validate, DB: db}
}

type service struct {
	DB              *gorm.DB
	config          *appconf.Config
	SubmissionsRepo repository.SubmissionsRepository
	validate        *validator.Validate
}

func (s service) Create(
	ctx context.Context, req *domain.SubmissionRequest,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if err := s.validate.Struct(req); err != nil {
		return exception.InvalidArgument(err)
	}

	var answers []domain.SubmissionResult
	for _, answer := range req.Answers {
		answers = append(answers, domain.SubmissionResult{
			UserId:     req.UserId,
			QuestionId: answer.QuestionId,
			Answer:     answer.Answer,
		})
	}
	jsonData, err := json.Marshal(answers)
	if err != nil {
		return exception.Internal("error marshalling old value", err)
	}
	body := &domain.Submissions{
		UserId:  req.UserId,
		Answers: jsonData,
	}
	body.DeclareRisk(answers)
	if err := s.SubmissionsRepo.Create(ctx, tx, body); err != nil {
		return exception.Internal("error inserting submissions", err)
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
	err = s.SubmissionsRepo.Delete(ctx, tx, idInt)
	if err != nil {
		return exception.Internal("error deleting submissions", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Find(ctx context.Context, limit, page string) (*FindResponse, *exception.Exception) {
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
	result, pagination, err := s.SubmissionsRepo.Find(ctx, tx, limitInt, pageInt)
	if err != nil {
		return nil, exception.Internal("error getting submissions", err)
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

func (s service) FindByUser(ctx context.Context, limit, page, userid string) (
	*FindByUserResponse, *exception.Exception,
) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(userid)
	if err != nil {
		return nil, exception.PermissionDenied("Input of id must be integer")
	}
	var limitInt, pageInt int
	limitInt, err = strconv.Atoi(limit)
	if err != nil {
		return nil, exception.PermissionDenied("Input of limit must be integer")
	}
	pageInt, err = strconv.Atoi(page)
	if err != nil {
		return nil, exception.PermissionDenied("Input of page must be integer")
	}

	result, pagination, err := s.SubmissionsRepo.FindByUser(ctx, tx, limitInt, pageInt, idInt)
	if err != nil {
		return nil, exception.Internal("error getting submissions", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	finalResponse := &FindByUserResponse{
		UserId:     idInt,
		Pagination: *pagination,
		Data:       *result,
	}
	return finalResponse, nil
}

func (s service) Detail(ctx context.Context, id string) (*domain.Submissions, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, exception.PermissionDenied("Input of id must be integer")
	}
	result, err := s.SubmissionsRepo.Detail(ctx, tx, idInt)
	if err != nil {
		return nil, exception.Internal("error getting detail submissions", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
