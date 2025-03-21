package user

import (
	"context"
	v1 "scaffold/api/user/v1"
	userLogic "scaffold/internal/logic/user"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	in := &userLogic.CreateIn{
		Username: req.Username,
		Email:    req.Email,
	}
	userId, err := c.userLogic.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return &v1.CreateRes{UserId: userId}, nil
}
