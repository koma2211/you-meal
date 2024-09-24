package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initCategoryHandler(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		categories.GET("/burgers/:page", h.getBurgersByPage())
	}
}

func (h *Handler) getBurgersByPage() gin.HandlerFunc {
	const pageParam = "page"
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param(pageParam))
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