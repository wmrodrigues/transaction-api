package handler

import (
	"net/http"
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
	user, err := uh.userService.GetById(ctx, id)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, user)
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest dto.User
	if err := common.DecodeBodyToObject(w, r, &userRequest); err != nil {
		return
	}
	var user domain.User
	user = userRequest.ToDomain()
	if err := uh.userService.Create(r.Context(), &user); err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	common.SendCreatedResponse(w, user)
}

func (uh *UserHandler) ValidateCredentials(w http.ResponseWriter, r *http.Request) {
	var userRequest dto.User
	if err := common.DecodeBodyToObject(w, r, &userRequest); err != nil {
		return
	}
	var user domain.User
	user = userRequest.ToDomain()
	if err := uh.userService.ValidateCredentials(r.Context(), user.Email, user.Password); err != nil {
		common.SendBadRequestResponse(w, err)
	}
	common.SendSuccessResponse(w, nil)
}
