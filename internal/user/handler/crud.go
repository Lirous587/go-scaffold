package handler

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/reskit/response"
)

func (h *HttpHandler) GetProfile(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := DomainUserToResponse(user)
	response.Success(ctx, res)
}

func (h *HttpHandler) UpdateProfile(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
	}

	req := new(UserProfileUpdateRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	updates := HTTPUserUpdateToDomain(req)
	user, err := h.userService.UpdateUserProfile(userID, updates)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := DomainUserToResponse(user)
	response.Success(ctx, res)
}
