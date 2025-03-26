package mock2

import (
	"context"
	v1 "scaffold/api/mock2/v1"
)

type IMock2V1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
}
