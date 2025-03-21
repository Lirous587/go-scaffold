package dao

import (
	"context"
	"scaffold/internal/model/entity"
	"scaffold/pkg/repository/db"
)

type User struct {
}

func New() *User {
	return &User{}
}

func (d *User) Insert(ctx context.Context, user *entity.User) (uint, error) {
	orm := db.DB()
	result := orm.WithContext(ctx).Create(user)
	return user.ID, result.Error
}
