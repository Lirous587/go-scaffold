package query

import (
	"context"
	"scaffold/internal/common/decorator"
	"scaffold/internal/feature/domain/user"
)

type LoginResponse struct {
}

type LoginByGithubQuery struct {
	Code string `json:"code" binding:"required"`
}

type LoginByGithubQueryHandler decorator.QueryHandler[LoginByGithubQuery, LoginResponse]

func NewLoginByGithubQueryHandler(userRepo user.Repository, metricsClient decorator.MetricsClient) LoginByGithubQueryHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}
	return decorator.ApplyQueryDecorators[LoginByGithubQuery, LoginResponse](
		loginByGithubQueryHandler{},
		metricsClient,
	)
}

type loginByGithubQueryHandler struct {
}

func (h loginByGithubQueryHandler) Handle(ctx context.Context, query LoginByGithubQuery) (LoginResponse, error) {
	return LoginResponse{}, nil
}
