package mock1

import (
	"context"
	"scaffold/api"
	v1 "scaffold/api/mock1/v1"
	"scaffold/response"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, api.Code) {
	return &v1.CreateRes{}, response.CodeSuccess
}
