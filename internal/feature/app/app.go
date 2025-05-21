package app

import (
	"scaffold/internal/feature/app/command"
	"scaffold/internal/feature/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	LoginByGithub command.LoginByGithubHandler
}

type Queries struct {
	UserJWTByGithub query.LoginByGithubQueryHandler
}
