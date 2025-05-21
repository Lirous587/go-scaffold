package command

type LoginType string

const (
	GithubLogin LoginType = "github"
	EmailLogin  LoginType = "email"
)
