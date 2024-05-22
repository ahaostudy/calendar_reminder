package model

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email          string `gorm:"unique" json:"email"`
	PasswordHashed string `gorm:"type:varchar(255);not null" json:"-"`
}

func GetById(db *gorm.DB, ctx context.Context, id uint) (user *User, err error) {
	err = db.WithContext(ctx).Model(&User{}).Where("id = ?", id).First(&user).Error
	return
}

func GetByEmail(db *gorm.DB, ctx context.Context, email string) (user *User, err error) {
	err = db.WithContext(ctx).Model(&User{}).Where("email = ?", email).First(&user).Error
	return
}

func Create(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}
