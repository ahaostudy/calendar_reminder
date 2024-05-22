package router

import (
	"github.com/ahaostudy/calendar_reminder/middleware"
	"github.com/gin-gonic/gin"

	"github.com/ahaostudy/calendar_reminder/controller"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/register/", controller.Register)
		api.POST("/user/login/", controller.Login)

		api.Use(middleware.AuthMiddleware())
		api.GET("/user", controller.Get)
	}
}
