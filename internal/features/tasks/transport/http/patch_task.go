package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/AlexeyBobkovDev/golang-todoapp/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"         swaggertype:"string" example:"Walk a dog"`
	Description core_http_types.Nullable[string] `json:"description"   swaggertype:"string" example:"Walk a dog at 6:30 am"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"     swaggertype:"string" example:"false"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can not be NULL")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask godoc
//
//	@Summary		Patch task
//	@Description	Update an information about a task already existing in the system
//	@Description	### Logic of update fields (Three-state logic):
//	@Description	1. **Field is not given**: `description` is ignored. The value is not changed in DB
//	@Description	2. **Field is given explicitly**: `"description": "At 6:30 am walk a dog"`
//	@Description	3. **Field is given as `null`**: `"description": null` - clear the field in the DB (set to NULL)
//	@Description	Constraints: `title` and `completed` can not be set as `null`
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"ID of user to patch"
//	@Param			request	body		PatchTaskRequest					true	"PatchTask body (fields to update)"
//	@Success		200		{object}	PatchTaskResponse					"Task was updated successfully"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"Not Found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}
	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
