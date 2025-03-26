package v1

import "scaffold/pkg/apigen"

type CreateReq struct {
	apigen.Meta   `method:"post" path:"mock1/{id}" sm:"mock1" dc:"mock1描述" tags:"mock1"`
	Authorization string `in:"header" validate:"required" dc:"Bearer令牌" example:"Bearer eyJhbGciOiJS..." default:"Bearer fuck"`
	JsonMock      string `json:"json_mock" validate:"required"`
	QueryMock     string `in:"query" query:"query_mock" validate:"required"`
	UriMock       int    `in:"path" uri:"id" validate:"required,lt=100"`
}

type CreateRes struct {
}
