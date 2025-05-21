package command

import (
	"context"
	"scaffold/internal/common/decorator"
	"scaffold/internal/feature/domain/user"
)

type LoginByGithub struct {
	Code string `json:"code" binding:"required"`
}

type LoginByGithubHandler decorator.CommandHandler[LoginByGithub]

type loginByGithubHandler struct {
	userRepo user.Repository
}

func NewLoginByGithubHandler(userRepo user.Repository, metricsClient decorator.MetricsClient) LoginByGithubHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}
	return decorator.ApplyCommandDecorators[LoginByGithub](
		loginByGithubHandler{userRepo: userRepo},
		metricsClient,
	)
}

func (h loginByGithubHandler) Handle(ctx context.Context, cmd LoginByGithub) error {
	return nil
}
