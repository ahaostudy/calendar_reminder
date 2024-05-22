package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/conf"
	"github.com/ahaostudy/calendar_reminder/log"
	"github.com/ahaostudy/calendar_reminder/middleware"
)

func init() {
	log.InitLogger()
}

func main() {
	r := gin.Default()
	initMiddleware(r)

	if err := r.Run(conf.GetConf().Server.Address); err != nil {
		logrus.Fatal(err)
	}
}

func initMiddleware(e *gin.Engine) {
	e.Use(middleware.LoggerMiddleware())
}
