package user

import (
	"context"
	v1 "scaffold/api/user/v1"
)


func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	// err := c.userLogic.Delete(ctx, req.Id)
	// if err != nil {
		// return nil, err
	// }

	return &v1.UpdateRes{Msg: "删除成功"}, nil
}
