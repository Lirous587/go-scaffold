package user

import (
	"context"
	v1 "scaffold/api/user/v1"
)

func (c *ControllerV1) Test(ctx context.Context, req *v1.TestReq) (*v1.TestRes, error) {
	res := &v1.TestRes{
		Code:    1000,
		Message: "fuck msg",
		Data:    "fuck data",
	}
	return res, nil
}
