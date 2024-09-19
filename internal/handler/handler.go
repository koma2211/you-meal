package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/you-meal/internal/service"

	"github.com/rs/zerolog"
)

type Handler struct {
	services     *service.Service
	accessLogger zerolog.Logger
	env          string
}

func NewHandler(
	services *service.Service,
	accessLogger zerolog.Logger,
	env string,
) *Handler {
	return &Handler{
		services:     services,
		accessLogger: accessLogger,
		env:          env,
	}
}

func (h *Handler) Init() *gin.Engine {
	if h.env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(h.requestLoggerHTTP(), gin.Recovery())

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	api := router.Group("/api")
	{
		_ = api.Group("/v1")
	}

}
