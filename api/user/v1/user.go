package v1

type CreateReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Label    string `uri:"label" binding:"required"`
}

type CreateRes struct {
	UserId uint `json:"user_id"`
}

