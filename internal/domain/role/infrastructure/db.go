package infrastructure

import (
	"gorm.io/gorm"
	"scaffold/internal/domain/role/model"
)

type IDB interface {
	CreateRole(Role *model.Role) error
	DeleteRole(id uint) error
	UpdateRole(id uint, req *model.UpdateRoleReq) error
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) IDB {
	//var role model.UserRole
	//orm.AutoMigrate(&role)
	//var casbinRole model.CasbinRule
	//orm.AutoMigrate(&casbinRole)

	return &db{orm: orm}
}

func (d *db) CreateRole(Role *model.Role) error {
	return d.orm.Create(Role).Error
}

func (d *db) DeleteRole(id uint) error {
	return d.orm.Model(&model.Role{}).Delete("where id = ?", id).Error
}

func (d *db) UpdateRole(id uint, req *model.UpdateRoleReq) error {
	return d.orm.Model(&model.Role{}).Where("id = ?", id).Updates(req).Error
}
