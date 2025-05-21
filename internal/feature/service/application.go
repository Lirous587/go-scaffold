package service

import (
	"context"
	"scaffold/internal/common/decorator"
	"scaffold/internal/feature/adapters"
	"scaffold/internal/feature/app"
	"scaffold/internal/feature/app/command"
	"scaffold/internal/feature/app/query"
)

func NewApplication(ctx context.Context, metricsClient decorator.MetricsClient) (app.Application, func()) {
	application, clean := newApplication(ctx, metricsClient)
	return application, func() {
		clean()
	}
}

func newApplication(ctx context.Context, metricsClient decorator.MetricsClient) (app.Application, func()) {
	psqlRepository := adapters.NewPSQLRepository()

	application := app.Application{
		Commands: app.Commands{
			LoginByGithub: command.NewLoginByGithubHandler(psqlRepository, metricsClient),
		},
		Queries: app.Queries{
			UserJWTByGithub: query.NewLoginByGithubQueryHandler(psqlRepository, metricsClient),
		},
	}

	clean := func() {
		//_ =
	}

	return application, clean
}
