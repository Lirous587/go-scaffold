package handler

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"resty.dev/v3"

	"scaffold/internal/common/response"
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

func (h *HttpHandler) GithubAuth(ctx *gin.Context) {
	req := new(GithubAuthRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	// 1. 获取 GitHub 用户信息
	userInfo, err := h.getGithubUserInfo(req.Code)
	if err != nil {
		response.Error(ctx, err) // 直接传递已包装的领域错误
		return
	}

	// 2. 调用业务逻辑
	session, err := h.userService.AuthenticateWithOAuth("github", userInfo)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 3. 转换为响应格式
	res := Domain2TokenToAuthResponse(session)
	response.Success(ctx, res)
}

func (h *HttpHandler) RefreshToken(ctx *gin.Context) {
	req := new(RefreshTokenRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	payload := domain.JwtPayload{
		UserID:     req.UserID,
		RandomCode: req.RandomCode,
	}

	session, err := h.userService.RefreshUserToken(payload, req.RefreshToken)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	res := DomainSessionToRefreshResponse(session)
	response.Success(ctx, res)
}

func (h *HttpHandler) getUserID(ctx *gin.Context) (string, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "", domain.ErrTokenInvalid
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", domain.ErrTokenInvalid
	}
	return userIDStr, nil
}

// GitHub API 调用逻辑 - 返回包装好的领域错误
func (h *HttpHandler) getGithubUserInfo(code string) (*domain.OAuthUserInfo, error) {
	accessToken, err := h.getGithubAccessToken(code)
	if err != nil {
		return nil, domain.NewGitHubAPIError("get_access_token", err)
	}

	userInfo, err := h.fetchGithubUserInfo(accessToken)
	if err != nil {
		return nil, domain.NewGitHubAPIError("get_user_info", err)
	}

	return userInfo, nil
}

func (h *HttpHandler) getGithubAccessToken(code string) (string, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", domain.ErrInvalidOAuthCode.WithDetail("reason", "missing_credentials")
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
		return "", err // 这里的错误会在上层被包装
	}

	if result.AccessToken == "" {
		return "", domain.ErrInvalidOAuthCode.WithDetail("reason", "empty_access_token")
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
		Name:     githubUser.Name,
		Email:    githubUser.Email,
		Avatar:   githubUser.AvatarURL,
	}, nil
}
