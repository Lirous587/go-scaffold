package service

import (
	"go.uber.org/zap"
	"os"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/utils"
	"time"

	"github.com/pkg/errors"

	"scaffold/internal/user/domain"
)

type userService struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

var (
	githubClientID     string
	githubClientSecret string
)

func NewUserService(userRepo domain.UserRepository, tokenService domain.TokenService) domain.UserService {
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if githubClientID == "" || githubClientSecret == "" {
		panic("加载环境变量失败")
	}
	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *userService) AuthenticateWithOAuth(provider string, userInfo *domain.OAuthUserInfo) (*domain.User2Token, error) {
	// 1. 查找或创建用户
	user, isNewUser, err := s.findOrCreateUserByOAuth(provider, userInfo)
	if err != nil {
		return nil, err
	}

	// 2. 更新最后登录时间（如果不是新用户）
	if !isNewUser {
		if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
			// 这个错误不应该阻止登录流程，记录日志即可
			zap.L().Error("更新用户最后登录时间失败", zap.String("user_id", user.ID), zap.Error(err))
		}
	}

	// 3. 生成 Token
	payload := domain.JwtPayload{
		UserID:     user.ID,
		RandomCode: utils.GenRandomCodeForJWT(),
	}

	accessToken, err := s.tokenService.GenerateAccessToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &domain.User2Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) RefreshUserToken(payload domain.JwtPayload, refreshToken string) (*domain.User2Token, error) {
	//1 . 生成新的 access token
	accessToken, err := s.tokenService.RefreshAccessToken(payload, refreshToken)
	if err != nil {
		return nil, err
	}

	//2. 刷新refresh token的时间
	if err := s.tokenService.ResetRefreshTokenExpiry(payload); err != nil {
		return nil, err
	}

	return &domain.User2Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) GetUser(userID string) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateUserProfile(userID string, updates *domain.UserProfileUpdate) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 应用更新
	if updates.Name != nil {
		user.Name = *updates.Name
	}
	if updates.Username != nil {
		// 检查用户名是否已被使用
		if exists, err := s.userRepo.UsernameExists(*updates.Username); err != nil {
			return nil, err
		} else if exists {
			return nil, codes.ErrUserAlreadyExists
		}
		user.Username = *updates.Username
	}
	if updates.Avatar != nil {
		user.AvatarURL = *updates.Avatar
	}

	return s.userRepo.Update(user)
}

// 私有辅助方法
func (s *userService) findOrCreateUserByOAuth(provider string, userInfo *domain.OAuthUserInfo) (*domain.User, bool, error) {
	// 1. 先通过 OAuth ID 查找
	user, err := s.userRepo.FindByOAuthID(provider, userInfo.ID)
	if err == nil {
		// 找到用户，更新信息
		return user, false, nil
	}

	if !errors.Is(err, codes.ErrUserNotFound) {
		return nil, false, errors.WithStack(err)
	}

	// 2. 通过邮箱查找现有用户
	if userInfo.Email != "" {
		user, err = s.userRepo.FindByEmail(userInfo.Email)
		if err == nil {
			// 绑定 OAuth 到现有用户
			user, err = s.bindOAuthToUser(user, provider, userInfo)
			return user, false, err
		}

		if !errors.Is(err, codes.ErrUserNotFound) {
			return nil, false, errors.WithStack(err)
		}
	}

	// 3. 创建新用户
	user, err = s.createUserFromOAuth(provider, userInfo)
	return user, true, err
}

func (s *userService) createUserFromOAuth(provider string, userInfo *domain.OAuthUserInfo) (*domain.User, error) {
	// 生成唯一用户名
	username, err := s.userRepo.GenerateUniqueUsername(userInfo.Login)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		Username:      username,
		AvatarURL:     userInfo.Avatar,
		EmailVerified: true, // OAuth 用户邮箱已验证
		Status:        "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 设置 OAuth ID
	switch provider {
	case "github":
		user.GithubID = userInfo.ID
	case "google":
		user.GoogleID = userInfo.ID
	case "gitlab":
		user.GitlabID = userInfo.ID
	}

	return s.userRepo.Create(user)
}

func (s *userService) bindOAuthToUser(user *domain.User, provider string, userInfo *domain.OAuthUserInfo) (*domain.User, error) {
	// 设置 OAuth ID
	switch provider {
	case "github":
		user.GithubID = userInfo.ID
	case "google":
		user.GoogleID = userInfo.ID
	case "gitlab":
		user.GitlabID = userInfo.ID
	}

	// 更新头像等信息
	if userInfo.Avatar != "" {
		user.AvatarURL = userInfo.Avatar
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}
