package role

import (
	"github.com/pkg/errors"
	"scaffold/internal/role/infrastructure"
	"scaffold/internal/role/model"
)

type Service interface {
	CreateRole(req *model.CreateRoleReq) (*model.Role, error)
	DeleteRole(id int) error
	UpdateRole(id int, req *model.UpdateRoleReq) error
	AllRole() ([]model.Role, error)
}

type service struct {
	repo  infrastructure.RoleRepository
	cache infrastructure.RoleCache
}

func NewService(repo infrastructure.RoleRepository, cache infrastructure.RoleCache) Service {
	return &service{repo: repo, cache: cache}
}

func (s *service) CreateRole(req *model.CreateRoleReq) (*model.Role, error) {
	role, err := s.repo.CreateRole(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return role, nil
}

func (s *service) DeleteRole(id int) error {
	if err := s.repo.DeleteRole(id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) UpdateRole(id int, req *model.UpdateRoleReq) error {
	if err := s.repo.UpdateRole(id, req); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) AllRole() ([]model.Role, error) {
	roles, err := s.repo.AllRole()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return roles, nil
}
