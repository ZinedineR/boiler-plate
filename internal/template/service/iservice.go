package service

import (
	"context"
	"ms-batch/internal/template/domain"
	"ms-batch/pkg/db"
)

type Service interface {
	// CRUD operations for MailTemplate
	CreateMailTemplate(ctx context.Context, model *domain.MailTemplate) error
	FindMailTemplate(ctx context.Context, limit, page int64, column, name string) (
		*[]domain.MailTemplate, *db.MongoPaginate, error,
	)
	FindOneMailTemplate(ctx context.Context, id string) (*domain.MailTemplate, error)
	FindOneMailTemplateByName(ctx context.Context, name string) (*domain.MailTemplate, error)
	UpdateMailTemplate(ctx context.Context, id string, model *domain.MailTemplate) error
	DeleteMailTemplate(ctx context.Context, id string) error
}
