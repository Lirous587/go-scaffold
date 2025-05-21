package user

type LoginType string

const (
	GithubLogin LoginType = "github"
	EmailLogin  LoginType = "email"
)
