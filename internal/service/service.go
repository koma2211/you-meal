package service

import (
	"context"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	cacherepository "github.com/koma2211/you-meal/internal/repository/cache_repository"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Categorier interface {
	GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error)
	GetBurgersCount(ctx context.Context, page int) (int, error)
}

type Service struct {
	Categorier
}

func NewService(
	repo *repository.Repository,
	cacheRepo *cacherepository.CacheRepository,
	logger *logger.Logger,
) *Service {
	return &Service{
		Categorier: NewCategoryService(repo.Categorier, cacheRepo.Categorier, logger),
	}
}
