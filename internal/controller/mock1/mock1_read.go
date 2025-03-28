package mock1

import (
	"context"
	"scaffold/api"
	v1 "scaffold/api/mock1/v1"
	"scaffold/response"
)

func (c *ControllerV1) Read(ctx context.Context, req *v1.ReadReq) (*v1.ReadRes, api.Code) {
	return &v1.ReadRes{}, response.CodeValidationError
}
