package db

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"scaffold/internal/domain/admin/model"
)

type DB interface {
	HaveOne() (bool, error)
	FindByEmail(email string) (*model.Admin, error)
	Create(admin *model.Admin) error
}

type db struct {
	orm *gorm.DB
}

func NewDB(orm *gorm.DB) DB {
	var admin model.Admin
	orm.AutoMigrate(&admin)
	return &db{orm: orm}
}

func (db *db) HaveOne() (bool, error) {
	var admin model.Admin
	err := db.orm.First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.WithStack(err)
	}
	if &admin != nil {
		return true, nil
	}
	return false, nil
}

func (db *db) FindByEmail(email string) (*model.Admin, error) {
	var admin model.Admin
	err := db.orm.Model(admin).Where("email = ?", email).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return &admin, nil
}

func (db *db) Create(admin *model.Admin) error {
	return db.orm.Create(admin).Error
}
