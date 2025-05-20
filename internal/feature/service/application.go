package service

import (
	"context"
	"scaffold/internal/common/decorator"
	"scaffold/internal/feature/adapters"
	"scaffold/internal/feature/app"
	"scaffold/internal/feature/app/command"
)

func NewApplication(ctx context.Context, metricsClient decorator.MetricsClient) app.Application {
	psqlRepository := adapters.NewPSQLRepository()

	return app.Application{
		Commands: app.Commands{
			LoginWithType:  command.NewLoginWithTypeHandle(psqlRepository, metricsClient),
			LoginWithType2: command.NewLoginWithTypeHandle2(psqlRepository, metricsClient),
			//CancelTraining:       command.NewCancelTrainingHandler(hourRepository, logger, metricsClient),
			//ScheduleTraining:     command.NewScheduleTrainingHandler(hourRepository, logger, metricsClient),
			//MakeHoursAvailable:   command.NewMakeHoursAvailableHandler(hourRepository, logger, metricsClient),
			//MakeHoursUnavailable: command.NewMakeHoursUnavailableHandler(hourRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			//HourAvailability:      query.NewHourAvailabilityHandler(hourRepository, logger, metricsClient),
			//TrainerAvailableHours: query.NewAvailableHoursHandler(datesRepository, logger, metricsClient),
		},
	}
}
