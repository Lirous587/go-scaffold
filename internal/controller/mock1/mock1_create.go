package mock1

import (
	"context"
	"errors"
	v1 "scaffold/api/mock1/v1"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	mockErr := errors.New("this is a mock err")
	return &v1.CreateRes{}, mockErr
}
