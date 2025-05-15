package permission

import (
	"comment/internal/domain/permission/infrastructure"
	"comment/internal/domain/role/model"
	"comment/pkg/response"
)

type IService interface {
	CreatePermission(req *model.CreateRoleReq) *response.AppError
	DeletePermission(id uint) *response.AppError
	UpdatePermission(id uint, req *model.UpdateRoleReq) *response.AppError
}

type service struct {
	db    infrastructure.IDB
	cache infrastructure.ICache
}

func NewService(db infrastructure.IDB, cache infrastructure.IDB) IService {
	return &service{db: db, cache: cache}
}

func (s *service) CreatePermission(req *model.CreateRoleReq) *response.AppError {
	//TODO implement me
	panic("implement me")
}

func (s *service) DeletePermission(id uint) *response.AppError {
	//TODO implement me
	panic("implement me")
}

func (s *service) UpdatePermission(id uint, req *model.UpdateRoleReq) *response.AppError {
	//TODO implement me
	panic("implement me")
}
