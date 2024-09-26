package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/config"
)

func PostgresDBConn(config *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func SetupPostgresDBConfig(dbConfig config.Database) (*pgxpool.Config, error) {
	connConfig, err := pgxpool.ParseConfig(dbConfig.Source)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = dbConfig.MaxConns
	connConfig.MinConns = dbConfig.MinConns
	connConfig.MaxConnLifetime = dbConfig.MaxConnLifeTime
	connConfig.MaxConnIdleTime = dbConfig.MaxConnIdleTime
	connConfig.HealthCheckPeriod = dbConfig.HealthCheckPeriod
	connConfig.ConnConfig.ConnectTimeout = dbConfig.ConnectTimeout

	return connConfig, nil
}
