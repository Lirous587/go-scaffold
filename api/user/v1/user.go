package v1

import "scaffold/pkg/apigen"

type CreateReq struct {
	apigen.Meta `method:"post" path:"auth/register" sm:"用户注册" dc:"创建新用户账户" tags:"user"`

	Username string `json:"username" binding:"required" dc:"用户名" example:"johndoe" minLength:"3" maxLength:"50"`
	Email    string `json:"email" binding:"required,email" dc:"邮箱地址" format:"email" example:"john@example.com"`
	Password string `json:"password" binding:"required" dc:"密码" format:"password" minLength:"8" maxLength:"64" example:"P@ssw0rd123"`
	Age      int    `json:"age" dc:"年龄" min:"18" max:"120" example:"30"`
	UserType string `json:"user_type" dc:"用户类型" enum:"personal,business,admin" default:"personal"`
	Referral string `json:"referral" dc:"推荐码" pattern:"^[A-Z]{2}[0-9]{6}$" example:"AB123456"`
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
	apigen.Meta   `method:"put" path:"user/{id}" sm:"更新用户" dc:"更新用户" tags:"user" security:"jwt"`
	Authorization string `in:"header" v:"required" dc:"Bearer令牌" example:"Bearer eyJhbGciOiJS..." default:"Bearer fuck"`
	Id            int    `uri:"id" binding:"required" dc:"用户ID"`
	Name          string `json:"name" dc:"用户名称"`
}

type UpdateRes struct {
	Msg string `json:"msg"`
}
