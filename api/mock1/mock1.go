package mock1

import (
	"context"
	"scaffold/api"
	v1 "scaffold/api/mock1/v1"
)

type IMock1V1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, code api.Code)
	Read(ctx context.Context, req *v1.ReadReq) (res *v1.ReadRes, code api.Code)
}
