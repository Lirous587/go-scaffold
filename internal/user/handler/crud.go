package handler

import (
	"github.com/gin-gonic/gin"
	"scaffold/internal/common/response"
	"scaffold/internal/user/domain"
)

func (h *HttpHandler) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, domain.ErrRefreshTokenInvalid)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.Error(ctx, domain.ErrRefreshTokenInvalid)
		return
	}

	user, err := h.userService.GetUser(userIDStr)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := DomainUserToResponse(user)
	response.Success(ctx, res)
}

func (h *HttpHandler) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Error(ctx, domain.ErrRefreshTokenInvalid)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.Error(ctx, domain.ErrRefreshTokenInvalid)
		return
	}

	req := new(UserProfileUpdateRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	updates := HTTPUserUpdateToDomain(req)
	user, err := h.userService.UpdateUserProfile(userIDStr, updates)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := DomainUserToResponse(user)
	response.Success(ctx, res)
}
