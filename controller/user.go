package controller

import (
	"github.com/ahaostudy/calendar_reminder/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	RegisterRequest struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		PasswordConfirm string `json:"password_confirm" binding:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

func Register(c *gin.Context) {
	req := new(RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("invalid params: %v", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}

	user, token, err := service.Register(c.Request.Context(), req.Email, req.Password, req.PasswordConfirm)
	if err != nil {
		logrus.Errorf("register failed: %v", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(gin.H{
		"token": token,
		"user":  user,
	}))
}

func Login(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("invalid params: %v", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}

	user, token, err := service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(gin.H{
		"token": token,
		"user":  user,
	}))
}

func Get(c *gin.Context) {
	id := c.GetUint("user_id")
	user, err := service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
	}
	c.JSON(http.StatusOK, Success(user))
}
