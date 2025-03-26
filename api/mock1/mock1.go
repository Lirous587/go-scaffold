package mock1

import (
	"context"
	v1 "scaffold/api/mock1/v1"
)

type IMock1V1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
}
