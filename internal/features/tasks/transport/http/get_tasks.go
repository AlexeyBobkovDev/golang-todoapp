package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
//
//	@Summary		Get tasks
//	@Description	Get tasks with an optional pagination or/and filtration by author id
//	@Tags			tasks
//	@Produce		json
//	@Param			user_id	query		int									false	"Filter tasks by user_id"
//	@Param			limit	query		int									false	"Size of page with tasks"
//	@Param			offset	query		int									false	"Shift pages with tasks"
//	@Success		200		{object}	GetTasksResponse					"Tasks were returned successfully"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found""
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `user_id` query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
