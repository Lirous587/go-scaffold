package command

import (
	"context"
	"scaffold/internal/common/decorator"
	"scaffold/internal/feature/domain/user"
)

type LoginWithType struct {
	LoginType string
}

type LoginWithTypeHandler decorator.CommandHandler[LoginWithType]

type loginWithTypeHandler struct {
	userRepo user.WriteRepository
}

func NewLoginWithTypeHandle(userRepo user.WriteRepository, metricsClient decorator.MetricsClient) LoginWithTypeHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}
	return decorator.ApplyCommandDecorators[LoginWithType](
		loginWithTypeHandler{userRepo: userRepo},
		metricsClient,
	)
}

func (h loginWithTypeHandler) Handle(ctx context.Context, cmd LoginWithType) error {
	return nil
}
