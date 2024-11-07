package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koma2211/you-meal/internal/config"
	"github.com/koma2211/you-meal/internal/handler"
	"github.com/koma2211/you-meal/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
	logger     *logger.Logger
}

func SetupServer(
	handler *handler.Handler,
	conf *config.HTTPServer,
	db *pgxpool.Pool,
	logger *logger.Logger,
) *Server {
	logger.DebugLog.Debug().Msg("server initialize...")

	httpServer := &http.Server{
		Addr:           conf.Address,
		Handler:        handler.Init(conf),
		ReadTimeout:    conf.ReadTimeOut,
		WriteTimeout:   conf.WriteTimeOut,
		IdleTimeout:    conf.IdleTimeout,
		MaxHeaderBytes: conf.MaxHeaderBytes,
	}

	logger.DebugLog.Debug().Msg(fmt.Sprintf("server listen by port:%s...", conf.Address))

	return &Server{httpServer: httpServer, db: db, logger: logger}
}

func (s *Server) Run() error {
	// Close db connection...
	defer s.db.Close()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
