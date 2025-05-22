package user

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/server/response"
	"scaffold/internal/user/model"
)

type Controller struct {
	server          Service
	loginStrategies map[model.LoginType]loginStrategy
}

func NewController(serve Service) *Controller {
	ctrl := &Controller{
		server:          serve,
		loginStrategies: make(map[model.LoginType]loginStrategy),
	}

	ctrl.loginStrategies[model.GithubLogin] = &githubLoginStrategy{}
	// 如果将来有新的登录方式，例如：
	// const googleLogin loginT = "google"
	// ctrl.loginStrategies[googleLogin] = &GoogleLoginStrategy{}

	return ctrl
}

func (c *Controller) Login(ctx *gin.Context) {
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
	login(ctx *gin.Context, c *Controller)
}

type githubLoginStrategy struct {
}

func (s *githubLoginStrategy) githubLoginStrategy(ctx *gin.Context, ctrl *Controller) {
	s.login(ctx, ctrl)
}

func (s *githubLoginStrategy) login(ctx *gin.Context, ctrl *Controller) {
	req := new(model.GithubAuthReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, err := ctrl.server.GithubAuth(req.Code)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, res)
}

func (c *Controller) RefreshToken(ctx *gin.Context) {
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

	res, err := c.server.RefreshToken(req, refreshToken)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, res)
}
