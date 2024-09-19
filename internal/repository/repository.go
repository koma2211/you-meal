package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Repository struct {}

func NewRepository(
	db *pgxpool.Pool,
	logger *logger.Logger,
) *Repository {
	return &Repository{}
}