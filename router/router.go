package router

import (
	"github.com/ahaostudy/calendar_reminder/middleware/ginmw"
	"github.com/gin-gonic/gin"

	"github.com/ahaostudy/calendar_reminder/controller"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/register", controller.Register)
		api.POST("/user/login", controller.Login)

		api.Use(ginmw.AuthMiddleware())
		api.GET("/user", controller.GetUser)

		api.GET("/task/:id", controller.GetTask)
		api.GET("/task", controller.ListTask)
		api.POST("/task", controller.CreateTask)
		api.PUT("/task/:id", controller.UpdateTask)
		api.DELETE("/task/:id", controller.DeleteTask)
		api.GET("/task/date", controller.ListTaskByDate)
	}
}
