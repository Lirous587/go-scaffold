package mock1

import (
	"context"
	v1 "scaffold/api/mock1/v1"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	return &v1.CreateRes{}, nil
}
