package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/service"
)

func GetTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("user_id")
	task, err := service.GetTask(c.Request.Context(), id, userId)
	if err != nil {
		logrus.Error("get task error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(task))
}

func ListTask(c *gin.Context) {
	userId := c.GetUint("user_id")
	tasks, err := service.GetTaskList(c.Request.Context(), userId)
	if err != nil {
		logrus.Error("get task list error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(tasks))
}

type CreateTaskRequest struct {
	Title string `json:"title" binding:"required"`
	Time  int64  `json:"time" binding:"required"`
}

func CreateTask(c *gin.Context) {
	req := new(CreateTaskRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Error("invalid params:", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}
	userId := c.GetUint("user_id")
	task, err := service.CreateTask(c.Request.Context(), userId, req.Title, req.Time)
	if err != nil {
		logrus.Error("create task error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(task))
}

type UpdateTaskRequest struct {
	ID    string
	Title string `json:"title"`
	Time  int64  `json:"time"`
}

func UpdateTask(c *gin.Context) {
	req := &UpdateTaskRequest{ID: c.Param("id")}
	if err := c.ShouldBind(req); err != nil {
		logrus.Error("invalid params:", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}
	logrus.Info(req)
	userId := c.GetUint("user_id")
	task, err := service.UpdateTask(c.Request.Context(), req.ID, userId, req.Title, req.Time)
	if err != nil {
		logrus.Error("update task error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(task))
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userId := c.GetUint("user_id")
	err := service.DeleteTask(c.Request.Context(), id, userId)
	if err != nil {
		logrus.Error("delete task error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

func ListTaskByDate(c *gin.Context) {
	d := c.Query("d")
	userId := c.GetUint("user_id")

	date, err := service.ParseDate(d)
	if err != nil {
		logrus.Error("invalid params:", err)
		c.JSON(http.StatusOK, WithStatusCode(StatusCodeInvalidParams))
		return
	}

	tasks, err := service.GetTaskListByDate(c.Request.Context(), userId, date)
	if err != nil {
		logrus.Error("get task list by date error:", err)
		c.JSON(http.StatusOK, WithStatus(StatusCodeOperationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(tasks))
}
