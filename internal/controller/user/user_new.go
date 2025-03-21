package user

import (
	user "scaffold/api/user"
	userLogic "scaffold/internal/logic/user"
)

type ControllerV1 struct {
	userLogic *userLogic.User
}

func NewV1() user.IUserV1 {
	return &ControllerV1{
		userLogic: userLogic.New(),
	}
}
