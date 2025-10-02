package server

import (
	"scaffold/internal/common/reskit/codes"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserID(ctx *gin.Context) (int64, error) {
	uidStr, exist := ctx.Get(UserIDKey)
	if !exist {
		return 0, codes.ErrUserNotFound
	}

	userID, ok := uidStr.(int64)
	if !ok {
		return 0, codes.ErrUserIDInvalid
	}

	return userID, nil
}

// const TenantIDKey = "tenant_id"

// func GetTenantID(ctx *gin.Context) (int64, error) {
// 	tenantIDStr := ctx.Param(TenantIDKey)

// 	if tenantIDStr == "" {
// 		return 0, codes.ErrTenantNotFound
// 	}

// 	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
// 	if err != nil {
// 		return 0, codes.ErrTenantIDInvalid.WithCause(err)
// 	}

// 	return tenantID, nil
// }
