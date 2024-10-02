package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/entities"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Categorier interface {
	GetBurgersByPage(ctx context.Context, limit, offset int) ([]entities.Burger, error)
	GetBurgerIngredientsById(ctx context.Context, burgerId int) ([]entities.Ingredient, error)
	GetNumberOfPagesByBurgers(ctx context.Context, limit int) (int, error)
	CheckExistenceImage(ctx context.Context, burgerId int, fileName string) (bool, error)
}


type Repository struct {
	Categorier
}

func NewRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *Repository {
	return &Repository{
		Categorier: NewCategoryRepository(db, logger),
	}
}