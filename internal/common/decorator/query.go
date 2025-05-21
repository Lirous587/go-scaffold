package decorator

import (
	"context"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
