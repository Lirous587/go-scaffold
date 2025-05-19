package app

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CancelTraining command.CancelTrainingHandler
}

type Queries struct {
	HourAvailability query.HourAvailabilityHandler
}
