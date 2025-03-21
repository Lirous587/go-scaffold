package v1

type CreateReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateRes struct {
	UserId uint `json:"user_id"`
}
