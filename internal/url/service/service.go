package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/url/domain"
	"boiler-plate/internal/url/repository"
	"boiler-plate/internal/url/repository/redisser"
	"boiler-plate/pkg/exception"
	"boiler-plate/pkg/random"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.UsersRepository, db *gorm.DB,
	redis redisser.RedisClient,
	validate *validator.Validate,
) Service {
	return &service{
		config: config, UsersRepo: repo, validate: validate, DB: db,
		RedisClient: redis,
	}
}

type service struct {
	DB          *gorm.DB
	RedisClient redisser.RedisClient
	config      *appconf.Config
	UsersRepo   repository.UsersRepository
	validate    *validator.Validate
}

func (s service) Create(
	ctx context.Context, req *domain.URL,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	if err := s.validate.Struct(req); err != nil {
		return exception.InvalidArgument(err)
	}
	req.Shorturl = random.GenerateCodes(req.Longurl)
	err := s.UsersRepo.Create(ctx, tx, req)
	if err != nil {
		return exception.Internal("error inserting users", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return exception.Internal("error marshalling user data", err)
	}
	if err := s.RedisClient.HSet(ctx, "URL", req.Shorturl, string(jsonData)); err != nil {
		return exception.Internal("redis hset", err)
	}
	return nil
}

func (s service) Detail(ctx context.Context, shorturl string) (*domain.URL, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	var result *domain.URL
	resultCache, err := s.RedisClient.HGet(ctx, "URL", shorturl)
	if resultCache == "" {
		result, err = s.UsersRepo.Detail(ctx, tx, shorturl)
		if err != nil {
			return nil, exception.Internal("error getting detail users", err)
		}
		if result == nil {
			return nil, exception.NotFound("url expired")
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			return nil, exception.Internal("error marshalling user data", err)
		}
		if err := s.RedisClient.HSet(ctx, "URL", result.Shorturl, string(jsonData)); err != nil {
			return nil, exception.Internal("redis hset", err)
		}
	}
	converter := []byte(resultCache)
	_ = json.Unmarshal(converter, &result)
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
