package ports

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"scaffold/internal/common/server/response"
	"scaffold/internal/feature/app"
	"scaffold/internal/feature/app/command"
	"scaffold/internal/feature/app/query"
)

type HttpServer struct {
	app             app.Application
	loginStrategies map[command.LoginType]loginStrategy
}

func NewHttpServer(application app.Application) HttpServer {
	loginStrategy := make(map[command.LoginType]loginStrategy)
	loginStrategy[command.GithubLogin] = &githubLoginStrategy{}

	return HttpServer{
		app:             application,
		loginStrategies: loginStrategy,
	}
}

func (h *HttpServer) Login(ctx *gin.Context) {
	loginTypeQuery, exist := ctx.GetQuery("type")
	if !exist {
		response.ErrorParameterInvalid(ctx, errors.New("缺少 'type' 查询参数"))
		return
	}
	currentLoginType := command.LoginType(loginTypeQuery)
	strategy, found := h.loginStrategies[currentLoginType]
	if !found {
		response.ErrorParameterInvalid(ctx, errors.New("不支持的登录类型: "+string(currentLoginType)))
		return
	}
	strategy.login(ctx, h)
}

type loginStrategy interface {
	login(ctx *gin.Context, h *HttpServer)
}

type githubLoginStrategy struct {
}

func (s *githubLoginStrategy) githubLoginStrategy(ctx *gin.Context, h *HttpServer) {
	s.login(ctx, h)
}

func (s *githubLoginStrategy) login(ctx *gin.Context, h *HttpServer) {
	req := new(query.LoginByGithubQuery)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrorParameterInvalid(ctx, err)
		return
	}

	res, err := h.app.Queries.UserJWTByGithub.Handle(ctx, *req)

	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, res)
}

func (h *HttpServer) RefreshToken(ctx *gin.Context) {
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
