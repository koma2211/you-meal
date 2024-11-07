package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initMealHandler(api *gin.RouterGroup) {
	meals := api.Group("/meals")
	{
		meals.GET("/image", h.getImageMeal())
	}
}

func (h *Handler) getImageMeal() gin.HandlerFunc {
	const (
		imageQueryParam = "image_path"
	)
	return func(c *gin.Context) {
		imageFileName := c.Query(imageQueryParam)

		filePath := h.imagePath + imageFileName

		data, err := os.ReadFile(filePath)
		if err != nil {
			response(c, http.StatusInternalServerError, "File not found", nil)
			return
		}

		// Set the Content-Type header for JPEG
		c.Header("Content-Type", "image/jpeg") // Set content type for JPEG
		c.Header("Content-Disposition", "attachment; filename="+imageFileName)
		c.Data(http.StatusOK, "image/jpeg", data)
	}
}
