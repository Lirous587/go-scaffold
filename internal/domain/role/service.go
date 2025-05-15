package role

import (
	"scaffold/internal/domain/permission/infrastructure"
	"scaffold/internal/domain/role/model"
	"scaffold/pkg/response"
)

type IService interface {
	CreateRole(req *model.CreateRoleReq) *response.AppError
	DeleteRole(id uint) *response.AppError
	UpdateRole(id uint, req *model.UpdateRoleReq) *response.AppError
}

type service struct {
	db    infrastructure.IDB
	cache infrastructure.ICache
}

func NewService(db infrastructure.IDB, cache infrastructure.ICache) IService {
	return &service{db: db, cache: cache}
}

func (s *service) CreateRole(req *model.CreateRoleReq) *response.AppError {
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.db.CreateRole(&role); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) DeleteRole(id uint) *response.AppError {
	if err := s.db.DeleteRole(id); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}

func (s *service) UpdateRole(id uint, req *model.UpdateRoleReq) *response.AppError {
	if err := s.db.UpdateRole(id, req); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}
	return nil
}
