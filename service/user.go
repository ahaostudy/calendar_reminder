package service

import (
	"context"
	"errors"
	"fmt"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/ahaostudy/calendar_reminder/dal/mysql"
	"github.com/ahaostudy/calendar_reminder/model"
	"github.com/ahaostudy/calendar_reminder/utils/jwt"
	"github.com/ahaostudy/calendar_reminder/utils/sha256"
)

func Register(ctx context.Context, email, password, passwordConfirm string) (*model.User, string, error) {
	if password != passwordConfirm {
		return nil, "", errors.New("password must be the same as confirm password")
	}
	user := &model.User{
		Email:          email,
		PasswordHashed: sha256.Encrypt(password),
	}
	if err := model.CreateUser(mysql.DB, ctx, user); err != nil {
		var e *mysqldriver.MySQLError
		// MySQL Error 1062: Duplicate entry
		if errors.As(err, &e) && e.Number == 1062 {
			return nil, "", errors.New("email is already registered")
		}
		return nil, "", fmt.Errorf("user register failed: %w", err)
	}
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := model.GetUserByEmail(mysql.DB, ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("user does not exists")
		}
		return nil, "", err
	}
	passwordHashed := sha256.Encrypt(password)
	if user.PasswordHashed != passwordHashed {
		return nil, "", errors.New("incorrect user name or password")
	}
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func GetUser(ctx context.Context, id uint) (*model.User, error) {
	user, err := model.GetUserById(mysql.DB, ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exists")
		}
		return nil, err
	}
	return user, nil
}
