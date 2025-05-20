package decorator

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := zap.L().With(
		zap.String("command", handlerType),
		zap.String("command_body", fmt.Sprintf("%#v", cmd)),
	)

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.Error("Failed to execute command", zap.String("error", err.Error()))
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base QueryHandler[C, R]
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := zap.L().With(
		zap.String("query", generateActionName(cmd)),
		zap.String("query_body", fmt.Sprintf("%#v", cmd)),
	)

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.Error("Failed to execute command", zap.String("error", err.Error()))
		}
	}()

	return d.base.Handle(ctx, cmd)
}
