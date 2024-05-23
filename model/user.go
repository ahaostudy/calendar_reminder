package model

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email          string `gorm:"unique" json:"email"`
	PasswordHashed string `gorm:"type:varchar(127);not null" json:"-"`
}

func GetUserById(db *gorm.DB, ctx context.Context, id uint) (user *User, err error) {
	err = db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return
}

func GetUserByEmail(db *gorm.DB, ctx context.Context, email string) (user *User, err error) {
	err = db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return
}

func CreateUser(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}
