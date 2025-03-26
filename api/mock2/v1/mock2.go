package v1

import (
	"scaffold/api"
)

type CreateReq struct {
	api.Meta `method:"post" path:"mock2" tags:"mock2"`
	Mock     string `json:"mock" validate:"required"`
}

type CreateRes struct {
}
