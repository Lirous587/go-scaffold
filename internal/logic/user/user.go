package user

import (
	"context"
	dao "scaffold/internal/dao/user"
	"scaffold/internal/model/entity"
)

type User struct {
	dao dao.User
}

func New() *User {
	return &User{}
}

type CreateIn struct {
	Username string
	Email    string
}

func (u *User) Create(ctx context.Context, req *CreateIn) (uint, error) {
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		
	}
	return u.dao.Insert(ctx, user)
}
