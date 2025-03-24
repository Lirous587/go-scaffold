package v1

import "scaffold/pkg/apigen"

type CreateReq struct {
	apigen.Meta `method:"post" path:"userx" sm:"创建用户" dc:"通过提供的用户信息创建新用户"`
	Username    string `json:"username" binding:"required" dc:"用户名" example:"johndoe" minLength:"3" maxLength:"50"`
	Email       string `json:"email" binding:"required,email" dc:"电子邮箱地址" example:"john@example.com" format:"email"`
	Password    string `json:"password" binding:"required" dc:"用户密码" example:"P@ssw0rd" format:"password" minLength:"8" maxLength:"64"`
	Role        string `json:"role" dc:"用户角色" enum:"admin,user,guest" default:"user"`
	Age         int    `json:"age" dc:"用户年龄" min:"18" max:"120" example:"30"`
}

type CreateRes struct {
	UserId uint `json:"user_id"`
}

type DeleteReq struct {
	apigen.Meta `method:"delete" path:"user/{id}" sm:"删除用户" dc:"删除用户" tags:"user"`
	Id          int `uri:"id" binding:"required"`
}

type DeleteRes struct {
	Msg string `json:"msg"`
}

type UpdateReq struct {
	apigen.Meta `method:"put" path:"user/{id}" sm:"更新用户" dc:"更新用户" tags:"user"`
	Id          int `uri:"id" binding:"required"`
}

type UpdateRes struct {
	Msg string `json:"msg"`
}
