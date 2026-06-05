package users_transport_http

import (
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/response"
)

// DeleteUser godoc

// @Summary		Delete the user
// @Description	Delete an existing user in the system by his ID
// @Tags			users
// @Param			id	path	int	true	"ID of the user to delete"
// @Success		204	"User successfully deleted"
// @Failure		400	{object}	core_http_response.ErrorResponse	"Bad request"
// @Failure		404	{object}	core_http_response.ErrorResponse	"User not found"
// @Failure		500	{object}	core_http_response.ErrorResponse	"Internal Server Error"
// @Router			/users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
