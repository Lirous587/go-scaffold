package user

import (
	"context"
	v1 "scaffold/api/user/v1"
)


func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	err := c.userLogic.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteRes{Msg: "删除成功"}, nil
}
