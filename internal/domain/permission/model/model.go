package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Resource    string `gorm:"size:60;not null"` // 资源路径，如 /api/user
	Action      string `gorm:"size:5;not null"`  // 操作类型，如 GET、POST
	Description string `gorm:"size:60"`
	MenuID      int    `gorm:""`
	//Roles       []Role `gorm:"many2many:role_permission"`
}

type CreatePermissionReq struct {
	Resource    string `gorm:"size:60;not null"` // 资源路径，如 /api/user
	Action      string `gorm:"size:5;not null"`  // 操作类型，如 GET、POST
	Description string `gorm:"size:60"`
	//Roles       []Role `gorm:"many2many:role_permission"`
}

type UpdatePermissionReq struct {
	Resource    string `gorm:"size:60;not null"` // 资源路径，如 /api/user
	Action      string `gorm:"size:5;not null"`  // 操作类型，如 GET、POST
	Description string `gorm:"size:60"`
	//Roles       []Role `gorm:"many2many:role_permission"`
}
