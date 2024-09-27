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
			burgers.GET("/:page", h.getBurgersByPage())
			burgers.GET("/pages-count", h.getNumberOfPagesByBurgers())
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