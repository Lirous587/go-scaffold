package adapters

import (
	"github.com/aarondl/null/v8"
	"scaffold/internal/common/orm"
	"scaffold/internal/user/domain"
)

func domainUserToORM(user *domain.User) *orm.User {
	if user == nil {
		return nil
	}

	ormUser := &orm.User{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}

	if user.PasswordHash != "" {
		ormUser.PasswordHash = null.StringFrom(user.PasswordHash)
	}

	if user.GithubID != "" {
		ormUser.GithubID = null.StringFrom(user.GithubID)
	}

	return ormUser
}

func ormUserToDomain(ormUser *orm.User) *domain.User {
	if ormUser == nil {
		return nil
	}

	user := &domain.User{
		ID:          ormUser.ID,
		Email:       ormUser.Email,
		Nickname:    ormUser.Nickname,
		CreatedAt:   ormUser.CreatedAt,
		UpdatedAt:   ormUser.UpdatedAt,
		LastLoginAt: ormUser.LastLoginAt,
	}

	if ormUser.PasswordHash.Valid {
		user.PasswordHash = ormUser.PasswordHash.String
	}

	if ormUser.GithubID.Valid {
		user.GithubID = ormUser.GithubID.String
	}

	return user
}
