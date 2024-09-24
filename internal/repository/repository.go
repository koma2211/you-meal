package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Categorier interface {
	GetBurgersByPage(ctx context.Context, limit, page int) ([]entities.Burger, error)
}

type Repository struct {
	Categorier
}

func NewRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *Repository {
	return &Repository{
		NewCategoryRepository(db, logger),
	}
}