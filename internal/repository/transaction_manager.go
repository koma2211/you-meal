package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	db *pgxpool.Pool
}

func NewTransactionManager(db *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	tx, err := tm.db.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tm *TransactionManager) Rollback(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}

func (tm *TransactionManager) Commit(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}
