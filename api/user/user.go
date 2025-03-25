package user

import (
	"context"
	v1 "scaffold/api/user/v1"
)

type IUserV1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
	Test(ctx context.Context, req *v1.TestReq) (res *v1.TestRes, err error)
}
