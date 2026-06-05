package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlexeyBobkovDev/golang-todoapp/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/logger"
	core_http_request "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/AlexeyBobkovDev/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Ivan Ivanovich"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

// PatchUser godoc
//
//	@Summary		Change the user
//	@Description	Change an information about the existing user by his ID
//	@Description	### Field update logic (three-state update semantics)
//	@Description	1. Field not provided: `phone_number` is ignored and not persisted to the database.
//	@Description	2. Explicit value provided: `"phone_number": "+79998887766"` updates the stored value.
//	@Description	3. Explicit null value: `"phone_number": null` clears the field (sets DB value to NULL).
//	@Description	Constraint: `full_name` must not be null.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"ID of the user to change"
//	@Param			request	body		PatchUserRequest					true	"PatchUser request body"
//	@Success		200		{object}	PatchUserResponse					"Successfully patched the user"
//	@Failure		400		{object}	core_http_response.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	core_http_response.ErrorResponse	"User Not Found"
//	@Failure		409		{object}	core_http_response.ErrorResponse	"Conflict"
//	@Failure		500		{object}	core_http_response.ErrorResponse	"Internal Server Error"
//	@Router			/users/{id} [patch]
func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can not be null")
		}

		fullNameLen := len([]rune(*r.FullName.Value))

		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))

			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must startswith `+` symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
