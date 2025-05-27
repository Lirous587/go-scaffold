package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"scaffold/internal/common/orm"
	"scaffold/internal/user/domain"
)

type PSQLUserRepository struct {
	db *sql.DB
}

func NewPSQLUserRepository(db *sql.DB) domain.UserRepository {
	return &PSQLUserRepository{db: db}
}

func (r *PSQLUserRepository) FindByID(userID string) (*domain.User, error) {
	ctx := context.Background()
	ormUser, err := orm.Users(orm.UserWhere.UserID.EQ(userID)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return ORMUserToDomain(ormUser), nil
}

func (r *PSQLUserRepository) FindByEmail(email string) (*domain.User, error) {
	ctx := context.Background()
	ormUser, err := orm.Users(orm.UserWhere.Email.EQ(email)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return ORMUserToDomain(ormUser), nil
}

func (r *PSQLUserRepository) Create(user *domain.User) (*domain.User, error) {
	ctx := context.Background()
	ormUser := DomainUserToORM(user)

	if err := ormUser.Insert(ctx, r.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return ORMUserToDomain(ormUser), nil
}

func (r *PSQLUserRepository) Update(user *domain.User) (*domain.User, error) {
	ctx := context.Background()
	ormUser := DomainUserToORM(user)

	_, err := ormUser.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return ORMUserToDomain(ormUser), nil
}

func (r *PSQLUserRepository) FindByOAuthID(provider, oauthID string) (*domain.User, error) {
	ctx := context.Background()
	var ormUser *orm.User
	var err error

	switch provider {
	case "github":
		ormUser, err = orm.Users(
			orm.UserWhere.GithubID.EQ(null.StringFrom(oauthID)),
		).One(ctx, r.db)
	case "google":
		ormUser, err = orm.Users(
			orm.UserWhere.GoogleID.EQ(null.StringFrom(oauthID)),
		).One(ctx, r.db)
	case "gitlab":
		ormUser, err = orm.Users(
			orm.UserWhere.GitlabID.EQ(null.StringFrom(oauthID)),
		).One(ctx, r.db)
	default:
		return nil, domain.ErrInvalidOAuthProvider
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return ORMUserToDomain(ormUser), nil
}

func (r *PSQLUserRepository) UpdateLastLogin(userID string) error {
	ctx := context.Background()
	ormUser, err := orm.Users(orm.UserWhere.UserID.EQ(userID)).One(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	ormUser.LastLoginAt = null.TimeFrom(time.Now())
	_, err = ormUser.Update(ctx, r.db, boil.Whitelist(orm.UserColumns.LastLoginAt))
	return err
}

func (r *PSQLUserRepository) EmailExists(email string) (bool, error) {
	ctx := context.Background()
	exists, err := orm.Users(orm.UserWhere.Email.EQ(email)).Exists(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}
	return exists, nil
}

func (r *PSQLUserRepository) UsernameExists(username string) (bool, error) {
	ctx := context.Background()
	exists, err := orm.Users(
		orm.UserWhere.Username.EQ(null.StringFrom(username)),
	).Exists(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}
	return exists, nil
}

func (r *PSQLUserRepository) GenerateUniqueUsername(preferred string) (string, error) {
	// 清理用户名：只保留字母数字和下划线
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return -1
	}, preferred)

	if cleaned == "" {
		cleaned = "user"
	}

	// 限制长度
	if len(cleaned) > 20 {
		cleaned = cleaned[:20]
	}

	// 检查是否已存在
	exists, err := r.UsernameExists(cleaned)
	if err != nil {
		return "", err
	}

	if !exists {
		return cleaned, nil
	}

	// 如果存在，添加数字后缀
	for i := 1; i <= 999; i++ {
		candidate := fmt.Sprintf("%s%d", cleaned, i)
		if len(candidate) > 30 {
			// 如果太长，截短基础名称
			maxBase := 30 - len(fmt.Sprintf("%d", i))
			candidate = fmt.Sprintf("%s%d", cleaned[:maxBase], i)
		}

		exists, err := r.UsernameExists(candidate)
		if err != nil {
			return "", err
		}

		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("unable to generate unique username")
}
