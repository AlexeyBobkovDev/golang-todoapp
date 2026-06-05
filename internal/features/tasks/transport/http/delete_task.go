package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/response"
)

// DeleteTask godoc
//
//	@Summary		Delete task
//	@Description	Delete the task from the system by ID
//	@Tags			tasks
//	@Param			id	path	int	true	"ID of task to delete"
//	@Success		204	"Task was deleted successfully"
//	@Failure		400	{object}	core_http_response.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		500	{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID Path value",
		)
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()
}
