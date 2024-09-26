package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/you-meal/internal/service"

	"github.com/rs/zerolog"
)

type Handler struct {
	services      *service.Service
	accessLogger  zerolog.Logger
	env           string
	limitCategory int
}

func NewHandler(
	services *service.Service,
	accessLogger zerolog.Logger,
	env string,
	limitCategory int,
) *Handler {
	return &Handler{
		services:     services,
		accessLogger: accessLogger,
		env:          env,
		limitCategory: limitCategory,
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
		v1 := api.Group("/v1")
		{
			h.initCategoryHandler(v1)
			h.initOrderHandler(v1)
		}
	}

}
