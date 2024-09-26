package handler

import "github.com/gin-gonic/gin"

func (h *Handler) initOrderHandler(api *gin.RouterGroup) {
	_ = api.Group("/orders")
	{}
}
