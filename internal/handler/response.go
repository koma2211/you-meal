package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koma2211/you-meal/internal/entities"
)

type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

func response(c *gin.Context, statusCode int, message string, data map[string]any) {
	if statusCode != 200 && statusCode != 400 && statusCode != 404 {
		switch message {
		case entities.ErrEmptyBurgers.Error():
			statusCode = http.StatusBadRequest
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
