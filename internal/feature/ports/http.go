package ports

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/pkg/response"
	"scaffold/internal/feature/app"
	"scaffold/internal/user/model"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}

}

func (h HttpServer) Login(ctx *gin.Context) {
	loginTypeQuery, exist := ctx.GetQuery("type")
	if !exist {
		response.ErrorParameterInvalid(ctx, errors.New("缺少 'type' 查询参数"))
		return
	}
	h.app.Commands.LoginWithType.han
	//currentLoginType := model.LoginType(loginTypeQuery)
	//strategy, found := c.loginStrategies[currentLoginType]
	//if !found {
	//	response.ErrorParameterInvalid(ctx, errors.New("不支持的登录类型: "+string(currentLoginType)))
	//	return
	//}
	//strategy.login(ctx, c)
}

func (h HttpServer) RefreshToken(ctx *gin.Context) {
	//refreshToken := ctx.GetHeader("refresh-token")
	//if refreshToken == "" {
	//	appError := response.NewAppError(response.CodeAuthFailed, errors.New("refresh-token请求头为空"))
	//	response.Error(ctx, appError)
	//	return
	//}
	//
	//req := new(model.RefreshTokenReq)
	//if err := ctx.ShouldBindJSON(req); err != nil {
	//	response.ErrorParameterInvalid(ctx, err)
	//	return
	//}
	//
	//res, err := c.server.RefreshToken(req, refreshToken)
	//if err != nil {
	//	response.Error(ctx, err)
	//	return
	//}
	//
	//response.Success(ctx, res)
}

type controller struct {
	loginStrategies map[model.LoginType]loginStrategy
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

	res, err := ctrl.server.Auth(req)
	if err != nil {
		response.Error(ctx, err)
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

	res, err := ctrl.server.Auth(req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, res)
}
