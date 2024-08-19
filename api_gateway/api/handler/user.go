package handler

import (
	_ "api_gateway/api/models"
	pbu "api_gateway/genproto/auth_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateUserProfileHandler godoc
// @Router       /user_service/v6/user/update/{id} [PUT]
// @Summary      Update user
// @Description  Update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "user_id"
// @Param        user body models.User false "user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) UpdateUserProfileHandler(ctx *gin.Context) {
	var (
		request  pbu.User
		response *pbu.User
		err      error
	)
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while reading by body of user_profile", http.StatusInternalServerError, err.Error())
		return
	}
	id := ctx.Param("id")
	request.Id = id
	response, err = h.services.UserService().UpdateUserProfile(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update user_profile", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, response)

}

// GetUserHandler      godoc
// @Router       /user_service/v6/user/{id} [GET]
// @Summary      Get user by id
// @Description  Get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "user_id"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) GetUserHandler(ctx *gin.Context) {
	var (
		request  pbu.PrimaryKeyUser
		err      error
		response *pbu.User
	)
	request.Id = ctx.Param("id")
	fmt.Println("++++++++", request.Id)
	response, err = h.services.UserService().GetUserProfile(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get user_profile", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, response)

}
