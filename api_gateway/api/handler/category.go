package handler

import (
	_ "api_gateway/api/docs"
	pbb "api_gateway/genproto/budgeting_service"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateCategoryHandler   godoc
// @Router       /budget_service/v3/category [POST]
// @Summary      Create a new category
// @Description  Create a new category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        category body models.CreateCategory true "category"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) CreateCategoryHandler(ctx *gin.Context) {
	var (
		request  pbb.CategoryRequest
		response *pbb.Category
		err      error
	)
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error while reading by body", http.StatusBadRequest, err.Error())
		return
	}
	response, err = h.services.CategoryService().CreateCategory(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while create category", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// UpdateCategoryHandler godoc
// @Router       /budget_service/v3/category/{id} [PUT]
// @Summary      Update category
// @Description  Update category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string true "category_id"
// @Param        category body models.UpdateCategory false "category"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) UpdateCategoryHandler(ctx *gin.Context) {
	var (
		request  pbb.Category
		response *pbb.Category
		id       string
		err      error
	)
	id = ctx.Param("id")
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error while reading by body", http.StatusBadRequest, err.Error())
		return
	}
	request.Id = id
	response, err = h.services.CategoryService().UpdateCategory(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update category", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// GetCategoryHandler      godoc
// @Router       /budget_service/v3/category/{id} [GET]
// @Summary      Get category by id
// @Description  Get category by id
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string true "category_id"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) GetCategoryHandler(ctx *gin.Context) {
	var (
		request  pbb.PrimaryKey
		response *pbb.Category
		err      error
	)
	request.Id = ctx.Param("id")
	response, err = h.services.CategoryService().GetCategory(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get category", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// DeleteCategoryHandler godoc
// @Router       /budget_service/v3/category/{id} [DELETE]
// @Summary      Delete category
// @Description  Delete category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id path string true "category_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) DeleteCategoryHandler(ctx *gin.Context) {
	var (
		err     error
		request pbb.PrimaryKey
	)
	request.Id = ctx.Param("id")
	_, err = h.services.CategoryService().DeleteCategory(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while delete category", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, "category successfully deleted!")

}

// GetAllCategoriesHandler godoc
// @Router       /budget_service/v3/categories [GET]
// @Summary      Get all category
// @Description  Get all category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.CategoriesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h Handler) GetAllCategoriesHandler(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	response, err := h.services.CategoryService().GetListCategories(context.Background(), &pbb.GetListRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error is while get all locations", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)
}
