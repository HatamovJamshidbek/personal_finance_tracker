package handler

import (
	_ "api_gateway/api/docs"
	pbb "api_gateway/genproto/budgeting_service"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05 -0700 MST"

// CreateGoalHandler   godoc
// @Router       /budget_service/v4/goal [POST]
// @Summary      Create a new goal
// @Description  Create a new goal
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        goal body models.CreateGoal true "goal"
// @Success      201  {object}  models.Goal
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) CreateGoalHandler(ctx *gin.Context) {
	var (
		request  pbb.GoalRequest
		response *pbb.Goal
		err      error
	)
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error while reading by body", http.StatusBadRequest, err.Error())
		return
	}
	request.Deadline += "T" + defaultTime + timeZoneOffset

	deadline, err := time.Parse(layout, request.GetDeadline())
	if err != nil {
		handleResponse(ctx, h.log, "start_time parse error", http.StatusBadRequest, err.Error())
		return
	}
	request.Deadline = deadline.String()
	response, err = h.services.GoalService().CreateGoal(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while create goal", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// UpdateGoalHandler godoc
// @Router       /budget_service/v4/goal/{id} [PUT]
// @Summary      Update goal
// @Description  Update goal
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        id path string true "goal_id"
// @Param        goal body models.UpdateGoal false "goal"
// @Success      200  {object}  models.Goal
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) UpdateGoalHandler(ctx *gin.Context) {
	var (
		request pbb.Goal
		id      string
		err     error
	)
	id = ctx.Param("id")
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error while reading by body", http.StatusBadRequest, err.Error())
		return
	}
	request.Id = id
	request.Deadline += "T" + defaultTime + timeZoneOffset

	deadline, err := time.Parse(layout, request.GetDeadline())
	if err != nil {
		handleResponse(ctx, h.log, "start_time parse error", http.StatusBadRequest, err.Error())
		return
	}
	request.Deadline = deadline.String()

	jsonData, err := json.Marshal(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update goal_progress is message broker that reading by body of goal", http.StatusInternalServerError, err.Error())
		return
	}
	err = h.rabbitMqProducer.ProduceMassage("goal_progress_updated", jsonData)
	if err != nil {
		handleResponse(ctx, h.log, "error is while update budget is message broker", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, "this success")

}

// GetGoalHandler      godoc
// @Router       /budget_service/v4/goal/{id} [GET]
// @Summary      Get goal by id
// @Description  Get goal by id
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        id path string true "goal_id"
// @Success      200  {object}  models.Goal
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) GetGoalHandler(ctx *gin.Context) {
	var (
		request  pbb.PrimaryKey
		response *pbb.Goal
		err      error
	)
	request.Id = ctx.Param("id")
	response, err = h.services.GoalService().GetGoal(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get goal", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "success", http.StatusCreated, response)

}

// DeleteGoalHandler godoc
// @Router       /budget_service/v4/goal/{id} [DELETE]
// @Summary      Delete goal
// @Description  Delete goal
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        id path string true "goal_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) DeleteGoalHandler(ctx *gin.Context) {
	var (
		err     error
		request pbb.PrimaryKey
	)
	request.Id = ctx.Param("id")
	_, err = h.services.GoalService().DeleteGoal(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while delete goal", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "Goal-Progress request successfully deleted", http.StatusCreated, "Goal-Progress  request successfully deleted")

}

// GetAllGoalsHandler godoc
// @Router       /budget_service/v4/goals [GET]
// @Summary      Get all goal
// @Description  Get all goal
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      200  {object}  models.GoalsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h Handler) GetAllGoalsHandler(c *gin.Context) {
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

	response, err := h.services.GoalService().GetListGoals(context.Background(), &pbb.GetListRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error is while get all goals", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)
}

// GenerateGoalProgressReportHandler godoc
// @Router       /budget_service/v4/goals/report-progress [GET]
// @Summary      Get all goal
// @Description  Get all goal
// @Tags         goal
// @Accept       json
// @Produce      json
// @Param        user_id query string false "user_id"
// @Param        start_time query string false "start_time"
// @Param        end_time query string false "end_time"
// @Success      200  {object}  models.GoalsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Security BearerAuth
func (h *Handler) GenerateGoalProgressReportHandler(ctx *gin.Context) {
	request := pbb.GoalProgressRequest{
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
	response, err := h.services.GoalService().GetGoalReportProgress(ctx, &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while get all goals-progress", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, response)

}
