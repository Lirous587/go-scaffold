package app

import "scaffold/internal/feature/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	LoginWithType command.LoginWithTypeHandler
}

type Queries struct {
	//HourAvailability query.HourAvailabilityHandler
}
