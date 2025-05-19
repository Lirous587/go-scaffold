package user

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"resty.dev/v3"
	"scaffold/internal/domain/user/infrastructure"
	"scaffold/internal/domain/user/model"
	"scaffold/pkg/config"
	"scaffold/pkg/jwt"
	"scaffold/pkg/response"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	Auth(req model.AuthReq) (*model.AuthRes, *response.AppError)
	RefreshToken(payload *model.JwtPayload, refreshToken string) (*model.RefreshTokenRes, *response.AppError)
}

type service struct {
	db    infrastructure.IDB
	cache infrastructure.ICache
}

func NewService(db infrastructure.IDB, cache infrastructure.ICache) IService {
	return &service{db: db, cache: cache}
}

func (s *service) genToken(payload *model.JwtPayload) (token string, err error) {
	jwtCfg := config.Cfg.JWT

	JWTTokenParams := jwt.JWTTokenParams{
		Payload:  *payload,
		Duration: time.Minute * time.Duration(jwtCfg.ExpireMinute),
		Secret:   []byte(jwtCfg.Secret),
	}

	token, err = jwt.GenToken[model.JwtPayload](&JWTTokenParams)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return
}

func (s *service) Auth(req model.AuthReq) (*model.AuthRes, *response.AppError) {
	switch req.(type) {
	case *model.EmailAuthReq:
		data := req.(*model.EmailAuthReq)
		return s.emailAuth(data.Email, data.Password)
	case *model.GithubAuthReq:
		data := req.(*model.GithubAuthReq)
		return s.githubAuth(data.Code)
	}
	return nil, nil
}

func (s *service) genTokenAndRefreshToken(payload *model.JwtPayload) (token, refreshToken string, err error) {
	jwtCfg := config.Cfg.JWT

	JWTTokenParams := jwt.JWTTokenParams{
		Payload:  *payload,
		Duration: time.Minute * time.Duration(jwtCfg.ExpireMinute),
		Secret:   []byte(jwtCfg.Secret),
	}

	token, err = jwt.GenToken[model.JwtPayload](&JWTTokenParams)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	refreshToken, err = s.cache.GenRefreshToken(payload)
	if err != nil {
		return "", "", response.NewAppError(response.CodeServerError, err)
	}

	return token, refreshToken, nil
}

func (s *service) emailAuth(email, password string) (*model.AuthRes, *response.AppError) {
	user, err := s.db.FindByEmail(email)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	if user == nil {
		return nil, response.NewAppError(response.CodeAuthFailed, err)
	}

	err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password))
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	payload := &model.JwtPayload{
		ID:        user.ID,
		LoginType: model.EmailLogin,
	}

	token, refreshToken, err := s.genTokenAndRefreshToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res := &model.AuthRes{
		Payload:      *payload,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (s *service) githubAuth(code string) (*model.AuthRes, *response.AppError) {
	client := resty.New()

	// 1. 获取 access_token
	oauthCfg := config.Cfg.OAuth.Github
	tokenURL := "https://github.com/login/oauth/access_token"
	params := map[string]string{
		"client_id":     oauthCfg.ClientID,
		"client_secret": oauthCfg.ClientSecret,
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
			return nil, response.NewAppError(response.CodeServerError, err)
		}
		return nil, response.NewAppError(response.CodeServerError, errors.New("获取github_token失败"))
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
		return nil, response.NewAppError(response.CodeServerError, err)
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
				return nil, response.NewAppError(response.CodeServerError, errors.WithStack(err))
			}
		} else {
			return nil, response.NewAppError(response.CodeServerError, err)
		}
	}

	payload := &model.JwtPayload{
		ID:        user.ID,
		LoginType: model.GithubLogin,
	}

	token, refreshToken, err := s.genTokenAndRefreshToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res := &model.AuthRes{
		Payload:      *payload,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (s *service) RefreshToken(payload *model.JwtPayload, refreshToken string) (res *model.RefreshTokenRes, appErr *response.AppError) {
	if err := s.cache.ValidateRefreshToken(payload, refreshToken); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, response.NewAppError(response.CodeRefreshInvalid, err)
		}
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	if err := s.cache.ResetRefreshTokenExpiry(payload); err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	newToken, err := s.genToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res = &model.RefreshTokenRes{
		Token: newToken,
	}
	return
}
