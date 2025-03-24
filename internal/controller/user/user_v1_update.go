package user

import (
	"context"
	v1 "scaffold/api/user/v1"
)

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	//  some logic

	return &v1.UpdateRes{Msg: "更新成功"}, nil
}
