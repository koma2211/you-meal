package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initCategoryHandler(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		categories.GET("/", h.getCategories())

		meals := categories.Group(":id/meals")
		{
			meals.GET("/", h.getMealsByCategoryID())
			meals.GET("/quantity", h.getNumberOfPagesByCategoryID())
		}
	}
}

// GetCategories Get all categories
//	@Summary		Get all categories
//	@Description	make a request for getting all categories
//	@ID				get-all-categories
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	Response
//	@Failure		500	{object}	Response
//	@Router			/categories/ [get]
func (h *Handler) getCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		menuCategories, err := h.services.Categorier.GetCategories(c.Request.Context())
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response(c, http.StatusOK, "success", map[string]any{"menuCategories": *menuCategories})
	}
}


// GetMealsByCategoryId Get meals by category id
//	@Summary		Get meals by category id
//	@Description	make a request for getting meals by category id with pagination
//	@ID				get-meals-by-category-id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Category ID"
//	@Param			page	query		int	true	"Page meal"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/categories/{id}/meals/ [get]
func (h *Handler) getMealsByCategoryID() gin.HandlerFunc {
	const (
		categoryIdQuery = "id"
		pageQuery       = "page"
	)
	return func(c *gin.Context) {
		categoryId, err := strconv.Atoi(c.Param(categoryIdQuery))
		if err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		page, err := strconv.Atoi(c.Query(pageQuery))
		if err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		meals, err := h.services.Categorier.GetMealsByCategoryID(c.Request.Context(), categoryId, h.limitCategory, page)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}


		response(c, http.StatusOK, "success", map[string]any{"meals": meals})
	}
}

// GetNumberOfPagesCategoryId Get number of pages category id
//	@Summary		Get number of pages category id
//	@Description	make a request for getting number of pages meals category id
//	@ID				get-number-of-pages-category-id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	Response
//	@Failure		400	{object}	Response
//	@Failure		500	{object}	Response
//	@Router			/categories/{id}/meals/quantity [get]
func (h *Handler) getNumberOfPagesByCategoryID() gin.HandlerFunc {
	const (
		categoryIdParam = "id"
	)
	return func(c *gin.Context) {
		categoryId, err := strconv.Atoi(c.Param(categoryIdParam))
		if err != nil {
			response(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		pagesCount, err := h.services.Categorier.GetMealPageCountByCategoryId(c.Request.Context(), categoryId, h.limitCategory)
		if err != nil {
			response(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response(c, http.StatusOK, "success", map[string]any{"pagesQuantity": pagesCount})
	}
}