package infrastructure

import (
	"scaffold/internal/user/model"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByGithubID(githubID int64) (*model.User, error)
	Register(u *model.User) (*model.User, error)
}
