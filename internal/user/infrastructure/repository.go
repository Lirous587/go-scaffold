package infrastructure

import (
	"scaffold/internal/user/model"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByGithubID(id string) (*model.User, error)
	Register(u *model.User) (*model.User, error)
}

//type db struct {
//	orm *gorm.DB
//}
//
//func NewDB(orm *gorm.DB) IDB {
//	//var user model.User
//	//orm.AutoMigrate(&user)
//	return &db{orm: orm}
//}
//
//func (db *db) FindByEmail(email string) (*model.User, error) {
//	var user model.User
//	err := db.orm.Model(user).Where("email = ?", email).First(&user).Error
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return &user, nil
//}
//
//func (db *db) FindByGithubID(id int) (*model.User, error) {
//	var user model.User
//	err := db.orm.Model(user).Where("github_id = ?", id).First(&user).Error
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return &user, nil
//}
//
//func (db *db) Create(user *model.User) error {
//	return db.orm.Create(user).Error
//}
