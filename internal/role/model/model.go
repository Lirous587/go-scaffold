package model

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description" json:"description,omitempty"`
}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required,max=64"`
	Description string `json:"description,omitempty" binding:"max=255"`
}

type UpdateRoleReq struct {
	Name        string `json:"name" binding:"required,max=64"`
	Description string `json:"description,omitempty" binding:"max=255"`
}
