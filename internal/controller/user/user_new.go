package user

import (
	user "scaffold/api/user"
	userLogic "scaffold/internal/logic/user"
	"scaffold/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ControllerV1 struct {
	userLogic *userLogic.User
}

func NewV1(r *gin.Engine) user.IUserV1 {
	ControllerV1 := &ControllerV1{
		userLogic: userLogic.New(),
	}

	group := r.Group("/v1")
	group.POST("/user/:label", middleware.AutoBind(ControllerV1.Create))
	
	return ControllerV1
}
