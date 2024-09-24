package service

import (
	"context"

	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Categorier interface {
	GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error)
}

type Service struct {
	Categorier
}

func NewService(
	repo *repository.Repository,
	logger *logger.Logger,
) *Service {
	return &Service{
		Categorier: NewCategoryService(repo.Categorier, logger),
	}
}
