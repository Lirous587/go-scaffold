package user

import (
	user "scaffold/api/user"
	userLogic "scaffold/internal/logic/user"
	"scaffold/pkg/apigen"

	"github.com/gin-gonic/gin"
)

type ControllerV1 struct {
	userLogic *userLogic.User
}

func NewV1(r *gin.Engine) user.IUserV1 {
	controller := &ControllerV1{
		userLogic: userLogic.New(),
	}

	// 使用自动API注册
	apigen.RegisterAPI(r, controller)

	return controller
}
