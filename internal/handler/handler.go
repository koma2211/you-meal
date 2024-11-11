package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/you-meal/docs"
	"github.com/koma2211/you-meal/internal/config"
	"github.com/koma2211/you-meal/internal/service"
	"github.com/rs/zerolog"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

type Handler struct {
	services      *service.Service
	accessLogger  zerolog.Logger
	env           string
	imagePath     string
	limitCategory int
}

func NewHandler(
	services *service.Service,
	accessLogger zerolog.Logger,
	env string,
	imagePath string,
	limitCategory int,
) *Handler {
	return &Handler{
		imagePath:     imagePath,
		services:      services,
		accessLogger:  accessLogger,
		env:           env,
		limitCategory: limitCategory,
	}
}

func (h *Handler) Init(cfg *config.HTTPServer) *gin.Engine {
	if h.env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(
		h.requestLoggerHTTP(),
		h.corsMiddleware(),
		gin.Recovery(),
	)

	// Init swagger
	docs.SwaggerInfo.Host = cfg.Address
	if cfg.Environment != config.EnvProd {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

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
			h.initMealHandler(v1)
		}
	}

}
