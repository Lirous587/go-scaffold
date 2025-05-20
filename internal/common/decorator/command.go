package decorator

import (
	"context"
	"fmt"
	"strings"
)

// ApplyCommandDecorators
// logger里面去调用 metrics metrics 里面去调用实际的handler
func ApplyCommandDecorators[H any](handler CommandHandler[H], metricsClient MetricsClient) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: metricsClient,
		},
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
