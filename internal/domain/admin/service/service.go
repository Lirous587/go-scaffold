package service

import (
	"github.com/redis/go-redis/v9"
	"scaffold/internal/domain/admin/infrastructure/cache"
	"scaffold/internal/domain/admin/infrastructure/db"
	"scaffold/internal/domain/admin/model"
	"scaffold/pkg/config"
	"scaffold/pkg/jwt"
	"scaffold/pkg/response"
	"scaffold/utils"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	IfInit() (bool, *response.AppError)
	Init(req *model.InitReq) *response.AppError
	Auth(email, password string) (*model.LoginRes, *response.AppError)
	RefreshToken(payload *model.JwtPayload, refreshToken string) (*model.RefreshTokenRes, *response.AppError)
}

type service struct {
	db    db.DB
	cache cache.Cache
}

func NewService(db db.DB, cache cache.Cache) Service {
	return &service{db: db, cache: cache}
}

func (s *service) IfInit() (bool, *response.AppError) {
	status, err := s.db.HaveOne()
	if err != nil {
		return status, response.NewAppError(response.CodeServerError, err)
	}
	return status, nil
}

func (s *service) Init(req *model.InitReq) *response.AppError {
	have, appErr := s.IfInit()
	if appErr != nil {
		return response.NewAppError(response.CodeServerError, appErr)
	}

	if have {
		return response.NewAppError(response.CodeAdminExist, errors.New("管理员已初始化"))
	}

	hashPassword, err := utils.EncryptPassword(req.Password)

	if err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	newAdmin := &model.Admin{
		Email:        req.Email,
		HashPassword: hashPassword,
	}

	if err = s.db.Create(newAdmin); err != nil {
		return response.NewAppError(response.CodeServerError, err)
	}

	return nil
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

func (s *service) Auth(email, password string) (*model.LoginRes, *response.AppError) {
	admin, err := s.db.FindByEmail(email)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	if admin == nil {
		return nil, response.NewAppError(response.CodeAuthFailed, err)
	}

	err = bcrypt.CompareHashAndPassword(admin.HashPassword, []byte(password))
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	payload := &model.JwtPayload{
		ID: admin.ID,
	}

	token, err := s.genToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	refreshToken, err := s.cache.GenRefreshToken(payload)
	if err != nil {
		return nil, response.NewAppError(response.CodeServerError, err)
	}

	res := &model.LoginRes{
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
