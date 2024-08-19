package handler

import (
	_ "api_gateway/api/docs"
	pbb "api_gateway/genproto/budgeting_service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateAccountHandler   godoc
// @Router       /budget_service/v1/account [POST]
// @Summary      Create a new account
// @Description  Create a new account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        account body models.CreateAccount true "account"
// @Success      201  {object}  models.Account
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) CreateAccountHandler(ctx *gin.Context) {
	var (
		request  pbb.AccountRequest
		response *pbb.Account
		err      error
	)
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error while reading by body", http.StatusBadRequest, err.Error())
		return
	}
	response, err = h.services.AccountService().CreateAccount(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while create account", http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("++++++++response", response)
	fmt.Println("id++++++++", response.Id)
	fmt.Println("created_at++++++++", response.Id)
	fmt.Println("++++++++type", response.Type)
	fmt.Println("++++++++user_id", response.UserId)
	fmt.Println("++++++++cur", response.Currency)
	fmt.Println("++++++++name", response.Name)
	fmt.Println("++++++++balance", response.Balance)
	fmt.Println("++++++++type request ---", request.Type)
	fmt.Println("++++++++user_id", request.UserId)
	fmt.Println("++++++++cur", request.Currency)
	fmt.Println("++++++++name", request.Name)
	fmt.Println("++++++++balance", request.Balance)
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// UpdateAccountHandler godoc
// @Router       /budget_service/v1/account/{id} [PUT]
// @Summary      Update account
// @Description  Update account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        id path string true "account_id"
// @Param        account body models.UpdateAccount false "account"
// @Success      200  {object}  models.Account
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) UpdateAccountHandler(ctx *gin.Context) {
	var (
		request  pbb.Account
		response *pbb.Account
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
	response, err = h.services.AccountService().UpdateAccount(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update account", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// GetAccountHandler      godoc
// @Router       /budget_service/v1/account/{id} [GET]
// @Summary      Get account by id
// @Description  Get account by id
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        id path string true "account_id"
// @Success      200  {object}  models.Account
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) GetAccountHandler(ctx *gin.Context) {
	var (
		request  pbb.PrimaryKey
		response *pbb.Account
		err      error
	)
	request.Id = ctx.Param("id")
	response, err = h.services.AccountService().GetAccount(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get account", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// DeleteAccountHandler godoc
// @Router       /budget_service/v1/account/{id} [DELETE]
// @Summary      Delete account
// @Description  Delete account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        id path string true "account_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) DeleteAccountHandler(ctx *gin.Context) {
	var (
		err     error
		request pbb.PrimaryKey
	)
	request.Id = ctx.Param("id")
	_, err = h.services.AccountService().DeleteAccount(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while delete account", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, "account successfully deleted!")

}

// GetAllAccountsHandler godoc
// @Router       /budget_service/v1/accounts [GET]
// @Summary      Get all account
// @Description  Get all account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.AccountsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h Handler) GetAllAccountsHandler(c *gin.Context) {
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

	response, err := h.services.AccountService().GetListAccounts(context.Background(), &pbb.GetListRequest{
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
