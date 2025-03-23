package entity

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;unique;not null"`
	Password string  `gorm:"size:255";not null`
}
