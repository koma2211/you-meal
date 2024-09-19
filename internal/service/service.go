package service

import (
	"github.com/koma2211/you-meal/internal/repository"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Service struct{}

func NewService(
	repo *repository.Repository,
	logger *logger.Logger,
) *Service {
	return &Service{}
}
