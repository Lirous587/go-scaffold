package user

import (
	"scaffold/api/user"
	userLogic "scaffold/internal/logic/user"
)

type ControllerV1 struct {
	userLogic *userLogic.User
}

func NewV1() user.IUserV1 {
	controller := &ControllerV1{
		userLogic: userLogic.New(),
	}
	return controller
}
