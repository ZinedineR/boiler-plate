package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"ms-batch/internal/template/domain"
	"ms-batch/internal/template/repository"
	"ms-batch/pkg/db"
	"ms-batch/pkg/errs"
)

// NewService creates new user service
func NewService(repo repository.Repository, validate *validator.Validate) Service {
	return &service{templateRepo: repo, validate: validate}
}

type service struct {
	templateRepo repository.Repository
	validate     *validator.Validate
}

func (s service) CreateMailTemplate(ctx context.Context, model *domain.MailTemplate) error {
	if err := s.validate.Struct(model); err != nil {
		return errs.Wrap(err)
	}
	if err := s.templateRepo.CreateMailTemplate(ctx, model); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (s service) FindMailTemplate(ctx context.Context, limit, page int64, column, name string) (
	*[]domain.MailTemplate, *db.MongoPaginate, error,
) {
	result, pagination, err := s.templateRepo.FindMailTemplate(ctx, limit, page, column, name)
	if err != nil {
		return nil, nil, errs.Wrap(err)
	}
	return result, pagination, nil
}

func (s service) FindOneMailTemplate(ctx context.Context, id string) (*domain.MailTemplate, error) {
	result, err := s.templateRepo.FindOneMailTemplate(ctx, id)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return result, nil
}

func (s service) FindOneMailTemplateByName(ctx context.Context, name string) (*domain.MailTemplate, error) {
	result, err := s.templateRepo.FindOneMailTemplateByName(ctx, name)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return result, nil
}

func (s service) UpdateMailTemplate(ctx context.Context, id string, model *domain.MailTemplate) error {
	if err := s.validate.Struct(model); err != nil {
		return errs.Wrap(err)
	}
	if err := s.templateRepo.UpdateMailTemplate(ctx, id, model); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (s service) DeleteMailTemplate(ctx context.Context, id string) error {
	if err := s.templateRepo.DeleteMailTemplate(ctx, id); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
