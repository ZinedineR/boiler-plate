package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/settings/domain"
	"boiler-plate/internal/settings/repository"
	"boiler-plate/pkg/errs"
	"context"
	"github.com/go-playground/validator/v10"
)

// NewService creates new user service
func NewService(config *appconf.Config, repo repository.Repository, validate *validator.Validate) Service {
	return &service{config: config, settingsRepo: repo, validate: validate}
}

type service struct {
	config       *appconf.Config
	settingsRepo repository.Repository
	validate     *validator.Validate
}

func (s service) FindSettings(ctx context.Context) (*domain.MainTable, error) {
	result, err := s.settingsRepo.FindSettings(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return result, nil
}
