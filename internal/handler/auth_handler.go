package handler

import (
	"fmt"
	"net/http"
	"transaction-api/internal/authentication"
	"transaction-api/internal/common"
	"transaction-api/internal/handler/dto"
	"transaction-api/internal/service"
)

type AuthHandler struct {
	userService  service.UserService
	tokenService authentication.TokenService
}

func NewAuthHandler(userService service.UserService, tokenService authentication.TokenService) *AuthHandler {
	return &AuthHandler{userService: userService, tokenService: tokenService}
}

func (h *AuthHandler) Tokens(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	err := common.DecodeBodyToObject(w, r, &req)
	if err != nil {
		common.SendBadRequestResponse(w, err)
		return
	}
	user, err := h.userService.ValidateCredentials(r.Context(), req.Email, req.Password)
	if err != nil {
		common.SendUnauthorizedResponse(w, fmt.Errorf("invalid credentials: %w", err))
		return
	}
	token, err := h.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		common.SendInternalServerErrorResponse(w, err)
		return
	}
	common.SendSuccessResponse(w, dto.AuthResponse{Token: token.Token, ExpiresIn: token.ExpiresAt, IssuedAt: token.IssuedAt})
}
