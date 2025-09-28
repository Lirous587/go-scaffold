package handler

import (
	"github.com/pkg/errors"
	"os"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/reskit/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"resty.dev/v3"

	"scaffold/internal/user/domain"
)

type HttpHandler struct {
	userService domain.UserService
}

func NewHttpHandler(userService domain.UserService) *HttpHandler {
	return &HttpHandler{
		userService: userService,
	}
}

// GithubAuth godoc
// @Summary      GitHub 授权登录
// @Description  使用 GitHub OAuth 登录，返回用户信息和令牌
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body handler.GithubAuthRequest true "GitHub 授权码"
// @Success      200 {object} response.successResponse{data=handler.AuthResponse} "登录成功"
// @Failure      400 {object} response.invalidParamsResponse "参数错误"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Router       /v1/user/auth/github [post]
func (h *HttpHandler) GithubAuth(ctx *gin.Context) {
	req := new(GithubAuthRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 1. 获取 GitHub 用户信息
	userInfo, err := h.getGithubUserInfo(req.Code)
	if err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// 2. 调用业务逻辑
	session, err := h.userService.AuthenticateWithOAuth("github", userInfo)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 3. 转换为响应格式
	response.Success(ctx, domain2TokenToAuthResponse(session))
}

func (h *HttpHandler) getRefreshToke(ctx *gin.Context) (string, error) {
	refreshToken := ctx.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		return "", codes.ErrRefreshTokenMissingInHeader
	}
	return refreshToken, nil
}

// RefreshToken godoc
// @Summary      刷新令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        X-Refresh-Token header string true "refresh_token刷新令牌"
// @Success      200 {object} response.successResponse{data=handler.RefreshTokenResponse} "刷新成功"
// @Failure      400 {object} response.errorResponse "参数错误"
// @Failure      401 {object} response.errorResponse "令牌无效或过期"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Router       /v1/user/refresh_token [post]
func (h *HttpHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := h.getRefreshToke(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	session, err := h.userService.RefreshUserToken(refreshToken)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := domainSessionToRefreshResponse(session)
	response.Success(ctx, res)
}

// GitHub API 调用逻辑 - 返回包装好的领域错误
func (h *HttpHandler) getGithubUserInfo(code string) (*domain.OAuthUserInfo, error) {
	accessToken, err := h.getGithubAccessToken(code)
	if err != nil {
		return nil, errors.WithStack(codes.ErrGitHubAPIError.WithSlug("get_access_token 获取失败").WithCause(err))
	}

	userInfo, err := h.fetchGithubUserInfo(accessToken)
	if err != nil {
		return nil, errors.WithStack(codes.ErrGitHubAPIError.WithSlug("get_user_info 获取失败").WithCause(err))
	}

	return userInfo, nil
}

func (h *HttpHandler) getGithubAccessToken(code string) (string, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", codes.ErrOAuthInvalidCode.WithDetail(map[string]any{
			"reason": "missing_credentials",
		})
	}

	client := resty.New()
	var result GithubAccessTokenResponse

	_, err := client.R().
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
		}).
		SetResult(&result).
		Post("https://github.com/login/oauth/access_token")

	if err != nil {
		return "", err
	}

	if result.AccessToken == "" {
		return "", codes.ErrOAuthInvalidCode.WithDetail(map[string]any{
			"reason": "empty_access_token",
		})
	}

	return result.AccessToken, nil
}

func (h *HttpHandler) fetchGithubUserInfo(accessToken string) (*domain.OAuthUserInfo, error) {
	client := resty.New()
	var githubUser GithubUser

	_, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Accept", "application/vnd.github+json").
		SetResult(&githubUser).
		Get("https://api.github.com/user")

	if err != nil {
		return nil, err // 这里的错误会在上层被包装
	}

	return &domain.OAuthUserInfo{
		Provider: "github",
		ID:       strconv.FormatInt(githubUser.ID, 10),
		Login:    githubUser.Login,
		Nickname: githubUser.Name,
		Email:    githubUser.Email,
		Avatar:   githubUser.AvatarURL,
	}, nil
}

func (h *HttpHandler) getUserID(ctx *gin.Context) (int64, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, codes.ErrTokenInvalid
	}

	userID64, ok := userID.(int64)
	if !ok {
		return 0, codes.ErrTokenInvalid
	}
	if userID64 == 0 {
		return 0, errors.New("无效的id")
	}
	return userID64, nil
}

// ValidateAuth godoc
// @Summary      校验令牌
// @Description  校验当前访问令牌是否有效
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.successResponse "令牌有效"
// @Failure      401 {object} response.errorResponse "令牌无效或过期"
// @Router       /v1/user/auth [post]
func (h *HttpHandler) ValidateAuth(ctx *gin.Context) {
	response.Success(ctx)
}

// GetProfile godoc
// @Summary      获取用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.successResponse{data=handler.UserResponse} "获取成功"
// @Failure      401 {object} response.errorResponse "未授权"
// @Failure      500 {object} response.errorResponse "服务器错误"
// @Router       /v1/user/profile [get]
func (h *HttpHandler) GetProfile(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		response.Error(ctx, err)
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, domainUserToResponse(user))
}
