package v1

import "scaffold/pkg/apigen"

type CreateReq struct {
	apigen.Meta `method:"post" path:"mock2" tags:"mock2"`
	Mock        string `json:"mock" validate:"required"`
}

type CreateRes struct {
}
