package handler

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

func response(c *gin.Context, statusCode int, message string, data map[string]any) {
	if statusCode != 200 && statusCode != 400 && statusCode != 404 {
		switch message {
		}
	}

	var response Response

	if statusCode != 200 {
		response = Response{
			Success: false,
			Message: message,
		}
	} else {
		response = Response{
			Success: true,
			Message: message,
			Data:    data,
		}
	}

	c.Header("Content-Type", "application/json")

	if statusCode != 200 {
		c.AbortWithStatusJSON(statusCode, response)
		return
	}

	c.JSON(statusCode, response)
}
