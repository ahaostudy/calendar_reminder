package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/service"
)

type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

func Register(c *gin.Context) {
	req := new(RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Error("invalid params:", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}

	user, token, err := service.Register(c.Request.Context(), req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		logrus.Error("register failed:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(gin.H{
		"token": token,
		"user":  user,
	}))
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Error("invalid params:", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}

	user, token, err := service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logrus.Error("login failed:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(gin.H{
		"token": token,
		"user":  user,
	}))
}

func GetUser(c *gin.Context) {
	id := c.GetUint("user_id")
	user, err := service.GetUser(c.Request.Context(), id)
	if err != nil {
		logrus.Error("get user failed:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(user))
}
