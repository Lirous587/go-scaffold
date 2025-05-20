package decorator

import (
	"context"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(action, status string, value int)
	ObserveDuration(action, status string, seconds float64)
}

type commandMetricsDecorator[C any] struct {
	base   CommandHandler[C]
	client MetricsClient
}

func (d commandMetricsDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	start := time.Now()

	actionName := strings.ToLower(generateActionName(cmd))

	defer func() {
		end := time.Since(start)
		status := "success"
		if err != nil {
			status = "failure"
		}
		d.client.Inc(actionName, status, 1)
		d.client.ObserveDuration(actionName, status, float64(end.Milliseconds()))
	}()

	return d.base.Handle(ctx, cmd)
}

type queryMetricsDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (d queryMetricsDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	start := time.Now()

	actionName := strings.ToLower(generateActionName(query))

	defer func() {
		end := time.Since(start)
		status := "success"
		if err != nil {
			status = "failure"
		}
		d.client.Inc(actionName, status, 1)
		d.client.ObserveDuration(actionName, status, end.Seconds())
	}()

	return d.base.Handle(ctx, query)
}
