package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/you-meal/internal/entities"
)

func (h *Handler) initOrderHandler(api *gin.RouterGroup) {
	orders := api.Group("/orders")
	{
		orders.POST("/", h.placeAnOrder())
	}
}

func (h *Handler) placeAnOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entities.Client

		if err := c.BindJSON(&req); err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		err := h.services.Customer.AddOrder(c.Request.Context(), req)

		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response(c, http.StatusOK, "success", nil)
	}
}
