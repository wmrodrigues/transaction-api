package handler

import (
	"fmt"
	"net/http"
	"transaction-api/internal/authentication"
	"transaction-api/internal/common"
	"transaction-api/internal/domain"
	"transaction-api/internal/handler/dto"
	"transaction-api/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	if claims.UserID != id {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: user id mismatch"))
		return
	}
	user, err := uh.userService.GetById(ctx, id)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, dto.UserToResponse(*user))
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest dto.UserRequest
	if err := common.DecodeBodyToObject(w, r, &userRequest); err != nil {
		return
	}
	var user domain.User
	user = userRequest.ToDomain()
	if err := uh.userService.Create(r.Context(), &user); err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendCreatedResponse(w, dto.UserToResponse(user))
}

func (uh *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := authentication.ClaimsFromContext(r.Context())
	if !ok {
		common.SendUnauthorizedResponse(w, fmt.Errorf("unauthorized: missing claims"))
		return
	}
	claimsMap := map[string]string{"user_id": claims.UserID, "email": claims.Email}
	common.SendSuccessResponse(w, claimsMap)
}
