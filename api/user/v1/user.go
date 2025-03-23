package v1

import "scaffold/api"

type CreateReq struct {
	api.Meta `method:"post" path:"v1/user" sm:"创建用户" dc:"通过提供的用户信息创建新用户" tags:"user"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateRes struct {
	UserId uint `json:"user_id"`
}

type DeleteReq struct {
	api.Meta `method:"delete" path:"v1/user/{id}" sm:"删除用户" dc:"删除用户" tags:"user"`
	Id       int `uri:"id" binding:"required"`
}

type DeleteRes struct {
	Msg string `json:"msg"`
}

type UpdateReq struct {
	api.Meta `method:"put" path:"v1/user/{id}" sm:"更新用户" dc:"更新用户" tags:"user"`
	Id       int `uri:"id" binding:"required"`
}

type UpdateRes struct {
	Msg string `json:"msg"`
}
