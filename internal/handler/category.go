package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initCategoryHandler(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		burgers := categories.Group("/burgers")
		{
			burgers.GET("/", h.getBurgersByPage())
			burgers.GET("/pages-count", h.getNumberOfPagesByBurgers())
			burgers.GET("/:burger_id", h.getImageBurgerById())
		}
	}
}

func (h *Handler) getBurgersByPage() gin.HandlerFunc {
	const pageQuery = "page"
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query(pageQuery))
		if err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		burgers, err := h.services.GetBurgersByPage(c.Request.Context(), h.limitCategory, page)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response(c, http.StatusOK, "success", map[string]any{"burgers": burgers})
	}
}

func (h *Handler) getNumberOfPagesByBurgers() gin.HandlerFunc {
	return func(c *gin.Context) {
		burgersPage, err := h.services.GetNumberOfPagesByBurgers(c.Request.Context(), h.limitCategory)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response(c, http.StatusOK, "success", map[string]any{"burgersPage": burgersPage})
	}
}

func (h *Handler) getImageBurgerById() gin.HandlerFunc {
	const (
		burgerIdParam   = "burger_id"
		imageQueryParam = "image_path"
	)
	return func(c *gin.Context) {
		burgerId, err := strconv.Atoi(c.Param(burgerIdParam))
		if err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		imagePath := c.Query(imageQueryParam)

		err = h.services.CheckExistenceImage(c.Request.Context(), burgerId, imagePath)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
}
