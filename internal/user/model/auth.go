package model

type LoginType string

const (
	GithubLogin LoginType = "github"
)

type JwtPayload struct {
	ID        int64     `json:"id"`
	LoginType LoginType `json:"login_type"`
}

type RefreshTokenReq = JwtPayload
type RefreshTokenRes struct {
	Token string `json:"token"`
}

type AuthRes struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	Payload      JwtPayload `json:"payload"`
}

type GithubAuthReq struct {
	Code string `json:"code" binding:"required"`
}
