package handler

import (
	_ "api_gateway/api/docs"
	pbb "api_gateway/genproto/budgeting_service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateTransactionHandler godoc
// @Router       /budget_service/v5/transaction [POST]
// @Summary      Create a new transaction
// @Description  Create a new transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        transaction body    models.CreateTransaction true "transaction"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) CreateTransactionHandler(ctx *gin.Context) {
	var request pbb.TransactionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleResponse(ctx, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	request.Date += "T" + defaultTime + timeZoneOffset

	date, err := time.Parse(layout, request.GetDate())
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}

	transactionMessage := pbb.TransactionRequest{
		UserId:      request.GetUserId(),
		AccountId:   request.GetAccountId(),
		CategoryId:  request.GetCategoryId(),
		Amount:      request.GetAmount(),
		Type:        request.GetType(),
		Description: request.GetDescription(),
		Date:        date.Format(time.RFC3339), // Ensure date is in RFC3339 format
	}

	// Marshal the request object to JSON
	jsonData, err := json.Marshal(&transactionMessage)
	if err != nil {
		handleResponse(ctx, h.log, "error marshalling request to JSON", http.StatusInternalServerError, err.Error())
		return
	}

	// Produce message to RabbitMQ
	if err := h.rabbitMqProducer.ProduceMassage("transaction_created", jsonData); err != nil {
		handleResponse(ctx, h.log, "error sending message to RabbitMQ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(ctx, h.log, "Transaction request successfully created", http.StatusCreated, "Transaction request successfully created")
}

// UpdateTransactionHandler godoc
// @Router       /budget_service/v5/transaction/{id} [PUT]
// @Summary      Update transaction
// @Description  Update transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction_id"
// @Param        transaction body models.UpdateTransaction false "transaction"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h *Handler) UpdateTransactionHandler(ctx *gin.Context) {
	var (
		request  pbb.Transaction
		response *pbb.Transaction
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
	request.Date += "T" + defaultTime + timeZoneOffset

	date, err := time.Parse(layout, request.GetDate())
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}
	request.Date = date.String()
	response, err = h.services.TransactionService().UpdateTransaction(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update transaction", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// GetTransactionHandler      godoc
// @Router       /budget_service/v5/transaction/{id} [GET]
// @Summary      Get transaction by id
// @Description  Get transaction by id
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction_id"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h *Handler) GetTransactionHandler(ctx *gin.Context) {
	var (
		request  pbb.PrimaryKey
		response *pbb.Transaction
		err      error
	)
	request.Id = ctx.Param("id")
	response, err = h.services.TransactionService().GetTransaction(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get transaction", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// DeleteTransactionHandler godoc
// @Router       /budget_service/v5/transaction/{id} [DELETE]
// @Summary      Delete transaction
// @Description  Delete transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        id path string true "transaction_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h *Handler) DeleteTransactionHandler(ctx *gin.Context) {
	var (
		err     error
		request pbb.PrimaryKey
	)
	request.Id = ctx.Param("id")
	_, err = h.services.TransactionService().DeleteTransaction(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while delete transaction", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, "transaction successfully deleted!")

}

// GetAllTransactionsHandler godoc
// @Router       /budget_service/v5/transactions [GET]
// @Summary      Get all transaction
// @Description  Get all transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.TransactionsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h Handler) GetAllTransactionsHandler(c *gin.Context) {
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

	response, err := h.services.TransactionService().GetListTransactions(context.Background(), &pbb.GetListRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
	})
	fmt.Println("response+++++", response)

	if err != nil {
		handleResponse(c, h.log, "error is while get all locations", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)
}

// GetUserIncomeHandler godoc
// @Router       /budget_service/v5/transactions/income [GET]
// @Summary      Get all transaction
// @Description  Get all transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        user_id query string false "user_id"
// @Param        start_time query string false "start_time"
// @Param        end_time query string false "end_time"
// @Success      200  {object}  models.GetUserMoneysResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h *Handler) GetUserIncomeHandler(ctx *gin.Context) {
	var (
		err      error
		response *pbb.GetUserMoneyResponse
	)

	userID := ctx.Query("user_id")
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if userID == "" || startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	startTimeStr += "T" + defaultTime + timeZoneOffset
	endTimeStr += "T" + defaultTime + timeZoneOffset

	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}

	request := pbb.GetUserMoneyRequest{
		UserId:    userID,
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	}

	// Fetch income data from service
	response, err = h.services.TransactionService().GetUserIncome(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error while getting user income", http.StatusInternalServerError, err.Error())
		return
	}

	// Send response
	handleResponse(ctx, h.log, "", http.StatusOK, response)
}

// GetUserSpendingHandler godoc
// @Router       /budget_service/v5/transactions/spend [GET]
// @Summary      Get all transaction
// @Description  Get all transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        user_id query string false "user_id"
// @Param        start_time query string false "start_time"
// @Param        end_time query string false "end_time"
// @Success      200  {object}  models.GetUserMoneysResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security 	 ApiKeyAuth
func (h *Handler) GetUserSpendingHandler(ctx *gin.Context) {
	var (
		err      error
		response *pbb.GetUserMoneyResponse
	)
	request := pbb.GetUserMoneyRequest{
		UserId:    ctx.Query("user_id"),
		StartTime: ctx.Query("start_time"),
		EndTime:   ctx.Query("end_time"),
	}
	request.StartTime += "T" + defaultTime + timeZoneOffset

	startTime, err := time.Parse(layout, request.GetStartTime())
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}
	request.EndTime += "T" + defaultTime + timeZoneOffset

	endTime, err := time.Parse(layout, request.GetEndTime())
	if err != nil {
		handleResponse(ctx, h.log, "date parse error", http.StatusBadRequest, err.Error())
		return
	}
	request.StartTime = startTime.String()
	request.EndTime = endTime.String()
	response, err = h.services.TransactionService().GetUserSpending(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get user spending money ", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, response)

}
