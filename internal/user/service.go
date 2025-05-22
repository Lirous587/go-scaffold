package user

import (
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"resty.dev/v3"
	"scaffold/internal/common/middleware/auth"
	"scaffold/internal/common/server/response"
	"scaffold/internal/user/infrastructure"
	"scaffold/internal/user/model"
)

type Service interface {
	GithubAuth(code string) (*model.AuthRes, error)
	RefreshToken(payload *model.JwtPayload, refreshToken string) (*model.RefreshTokenRes, error)
}

type service struct {
	db    infrastructure.UserRepository
	cache infrastructure.UserCache
}

var (
	githubClientID     string
	githubClientSecret string
)

func NewService(db infrastructure.UserRepository, cache infrastructure.UserCache) Service {
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if githubClientID == "" || githubClientSecret == "" {
		panic("加载环境变量失败")
	}
	return &service{db: db, cache: cache}
}

func (s *service) genTokenAndRefreshToken(payload *model.JwtPayload) (token, refreshToken string, err error) {
	token, err = auth.GenUserToken(payload)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	refreshToken, err = s.cache.GenRefreshToken(payload)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	return token, refreshToken, nil
}

func (s *service) GithubAuth(code string) (*model.AuthRes, error) {
	client := resty.New()

	// 1. 获取 access_token
	//oauthCfg := config.Cfg.OAuth.Github
	tokenURL := "https://github.com/login/oauth/access_token"

	params := map[string]string{
		"client_id":     githubClientID,
		"client_secret": githubClientSecret,
		"code":          code,
	}

	var tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	_, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetResult(&tokenRes).
		Post(tokenURL)

	if err != nil || tokenRes.AccessToken == "" {
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.New("获取github_token失败")
	}

	// 2. 用 access_token 获取用户信息
	var userInfo struct {
		GithubID int    `json:"id"`
		Name     string `json:"login"`
		Email    string `json:"email"`
	}
	_, err = client.R().
		SetHeader("Authorization", "Bearer "+tokenRes.AccessToken).
		SetResult(&userInfo).
		Get("https://api.github.com/user")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	//	3.是否入库
	user, err := s.db.FindByGithubID(userInfo.GithubID)
	if err != nil {
		// 没有就入库
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u := &model.User{
				Email:    userInfo.Email,
				GithubID: userInfo.GithubID,
			}

			if err := s.db.Create(u); err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}

	payload := &model.JwtPayload{
		ID:        user.ID,
		LoginType: model.GithubLogin,
	}

	token, refreshToken, err := s.genTokenAndRefreshToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := &model.AuthRes{
		Payload:      *payload,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (s *service) RefreshToken(payload *model.JwtPayload, refreshToken string) (*model.RefreshTokenRes, error) {
	if err := s.cache.ValidateRefreshToken(payload, refreshToken); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, response.NewAppError(response.CodeRefreshInvalid, err)
		}
		return nil, errors.WithStack(err)
	}

	if err := s.cache.ResetRefreshTokenExpiry(payload); err != nil {
		return nil, errors.WithStack(err)
	}

	newToken, err := auth.GenUserToken(payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := &model.RefreshTokenRes{
		Token: newToken,
	}
	return res, nil
}
