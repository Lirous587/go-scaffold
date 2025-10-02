package adapters

import (
	"database/sql"
	"fmt"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"scaffold/internal/common/reskit/codes"
	"time"

	"scaffold/internal/common/orm"
	"scaffold/internal/user/domain"
)

type UserPSQLRepository struct {
}

func NewUserPSQLRepository() domain.UserRepository {
	return &UserPSQLRepository{}
}

func (r *UserPSQLRepository) FindByID(id int64) (*domain.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.ID.EQ(id)).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return ormUserToDomain(ormUser), nil
}

func (r *UserPSQLRepository) FindByEmail(email string) (*domain.User, error) {
	ormUser, err := orm.Users(orm.UserWhere.Email.EQ(email)).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return ormUserToDomain(ormUser), nil
}

func (r *UserPSQLRepository) Create(user *domain.User) (*domain.User, error) {
	ormUser := domainUserToORM(user)

	if err := ormUser.InsertG(boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return ormUserToDomain(ormUser), nil
}

func (r *UserPSQLRepository) Update(user *domain.User) (*domain.User, error) {
	ormUser := domainUserToORM(user)

	_, err := ormUser.UpdateG(boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return ormUserToDomain(ormUser), nil
}

func (r *UserPSQLRepository) FindByOAuthID(provider, oauthID string) (*domain.User, error) {
	var ormUser *orm.User
	var err error

	switch provider {
	case "github":
		ormUser, err = orm.Users(
			orm.UserWhere.GithubID.EQ(null.StringFrom(oauthID)),
		).OneG()
	default:
		return nil, codes.ErrOAuthInvalidCode
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return ormUserToDomain(ormUser), nil
}

func (r *UserPSQLRepository) UpdateLastLogin(id int64) error {
	ormUser, err := orm.Users(orm.UserWhere.ID.EQ(id)).OneG()
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	ormUser.LastLoginAt = time.Now()
	_, err = ormUser.UpdateG(boil.Whitelist(orm.UserColumns.LastLoginAt))
	return err
}

func (r *UserPSQLRepository) EmailExists(email string) (bool, error) {
	exists, err := orm.Users(orm.UserWhere.Email.EQ(email)).ExistsG()
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}
	return exists, nil
}
