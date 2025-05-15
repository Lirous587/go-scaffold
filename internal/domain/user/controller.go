package user

import (
	"comment/internal/domain/user/model"
	"comment/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type IController interface {
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type controller struct {
	server          IService
	loginStrategies map[model.LoginType]loginStrategy
}

func NewController(svc IService) IController {
	ctrl := &controller{
		server:          svc,
		loginStrategies: make(map[model.LoginType]loginStrategy),
	}

	ctrl.loginStrategies[model.GithubLogin] = &githubLoginStrategy{}
	ctrl.loginStrategies[model.EmailLogin] = &emailLoginStrategy{}
	// 如果将来有新的登录方式，例如：
	// const googleLogin loginT = "google"
	// ctrl.loginStrategies[googleLogin] = &GoogleLoginStrategy{}

	return ctrl
}

func (c *controller) Login(ctx *gin.Context) {
	loginTypeQuery, exist := ctx.GetQuery("type")
	if !exist {
		response.ErrorParameterInvalid(ctx, errors.New("缺少 'type' 查询参数"))
		return
	}

	currentLoginType := model.LoginType(loginTypeQuery)
	strategy, found := c.loginStrategies[currentLoginType]
	if !found {
		response.ErrorParameterInvalid(ctx, errors.New("不支持的登录类型: "+string(currentLoginType)))
		return
	}
	strategy.login(ctx, c)
}

type loginStrategy interface {
	login(ctx *gin.Context, c *controller)
}

type githubLoginStrategy struct {
}

func (s *githubLoginStrategy) githubLoginStrategy(ctx *gin.Context, ctrl *controller) {
	s.login(ctx, ctrl)
}

func (s *githubLoginStrategy) login(ctx *gin.Context, ctrl *controller) {
	req := new(model.GithubAuthReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, appErr := ctrl.server.Auth(req)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}
	response.Success(ctx, res)
}

type emailLoginStrategy struct {
}

func (s *emailLoginStrategy) loginStrategy(ctx *gin.Context, ctrl *controller) {
	s.login(ctx, ctrl)
}

func (s *emailLoginStrategy) login(ctx *gin.Context, ctrl *controller) {
	req := new(model.EmailAuthReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, appErr := ctrl.server.Auth(req)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}
	response.Success(ctx, res)
}

func (c *controller) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("refresh-token")
	if refreshToken == "" {
		appError := response.NewAppError(response.CodeAuthFailed, errors.New("refresh-token请求头为空"))
		response.Error(ctx, appError)
		return
	}

	req := new(model.RefreshTokenReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, appErr := c.server.RefreshToken(req, refreshToken)
	if appErr != nil {
		response.Error(ctx, appErr)
		return
	}

	response.Success(ctx, res)
}
