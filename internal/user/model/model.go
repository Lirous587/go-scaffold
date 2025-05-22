package model

type User struct {
	ID       int
	Name     string
	Email    string
	GithubID string
	//Roles        []model.Role `gorm:"many2many:user_roles"`
}
