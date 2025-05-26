package infrastructure

import (
	"scaffold/internal/role/model"
)

type RoleRepository interface {
	CreateRole(req *model.CreateRoleReq) (*model.Role, error)
	DeleteRole(id int) error
	UpdateRole(id int, req *model.UpdateRoleReq) error
	AllRole() ([]model.Role, error)
}
