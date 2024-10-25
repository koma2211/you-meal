package main

import (
	"log"

	"github.com/koma2211/you-meal/internal/config"
	"github.com/koma2211/you-meal/internal/handler"
	"github.com/koma2211/you-meal/internal/repository"
	cacherepository "github.com/koma2211/you-meal/internal/repository/cache_repository"
	"github.com/koma2211/you-meal/internal/server"
	"github.com/koma2211/you-meal/internal/service"
	"github.com/koma2211/you-meal/pkg/cache/redis"
	"github.com/koma2211/you-meal/pkg/database/migration"
	"github.com/koma2211/you-meal/pkg/database/postgres"
	"github.com/koma2211/you-meal/pkg/logger"
	"github.com/koma2211/you-meal/pkg/validate"
)

func main() {
	cfg := config.MustLoad()

	accessLogger := logger.InitAccessLog(cfg.Logs)

	logger, err := logger.InitLogger(cfg.Logs)
	if err != nil {
		log.Fatalf("error when to init logger: %s", err.Error())
	}

	if err := migration.InitMigrationDB(cfg.MigrateSource, cfg.Database.Source); err != nil {
		log.Fatalf("error when to init migration: %s", err.Error())
	}

	dbConfig, err := postgres.SetupPostgresDBConfig(cfg.Database)
	if err != nil {
		log.Fatalf("eror when to setup postgres config: %s", err.Error())
	}

	db, err := postgres.PostgresDBConn(dbConfig)
	if err != nil {
		log.Fatalf("error when to connect postgres-db: %s", err.Error())
	}

	rdb, err := redis.RedisCacheMemoryConn(cfg.RedisSource)
	if err != nil {
		log.Fatalf("eror when  to connect redis: %s", err.Error())
	}

	cache := cacherepository.NewCacheRepository(rdb, logger, cfg.CacheCategoryTTL)

	valid := validate.NewValidation()

	repos := repository.NewRepository(db, logger)
	serv := service.NewService(repos, cache, logger, valid, cfg.RecievingTTL)
	handlers := handler.NewHandler(serv, accessLogger, cfg.Env, cfg.ImagePath, cfg.LimitCategory)

	server := server.SetupServer(handlers, &cfg.HTTPServer, db, logger)

	if err := server.Run(); err != nil {
		log.Fatalf("error when to run server: %s", err.Error())
	}
}
