package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

		imageFileName := c.Query(imageQueryParam)

		err = h.services.CheckExistenceImage(c.Request.Context(), burgerId, imageFileName)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		filePath := h.imagePath + imageFileName

		file, err := os.Open(filePath)
		if err != nil {
			response(c, http.StatusInternalServerError, fmt.Sprintf("Error opening file: %s", err.Error()), nil)
			return
		}

		defer file.Close()

		// Create a buffer to hold the multipart response
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create a form file part
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			response(c, http.StatusInternalServerError, fmt.Sprintf("Error creating form file: %s", err.Error()), nil)
			return
		}

		// Copy the file contents to the form file part
		if _, err := io.Copy(part, file); err != nil {
			response(c, http.StatusInternalServerError, fmt.Sprintf("Error copying file: %s", err.Error()), nil)
			return
		}

		// Close the writer to finalize the multipart data
		writer.Close()

		// Set the appropriate headers
		c.Header("Content-Type", writer.FormDataContentType())
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name()))

		// Send the multipart response
		c.Data(http.StatusOK, writer.FormDataContentType(), body.Bytes())
	}
}
