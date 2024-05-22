package main

import (
	"github.com/ahaostudy/calendar_reminder/dal/mysql"
	"github.com/ahaostudy/calendar_reminder/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/router"

	"github.com/ahaostudy/calendar_reminder/conf"
	"github.com/ahaostudy/calendar_reminder/middleware"

	"github.com/joho/godotenv"
)

func init() {
	log.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	r := gin.Default()

	// middleware
	initMiddleware(r)
	// router
	router.InitRouter(r)
	// mysql
	mysql.InitMySQL()

	if err := r.Run(conf.GetConf().Server.Address); err != nil {
		logrus.Fatal(err)
	}
}

func initMiddleware(r *gin.Engine) {
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.GlobalAuthMiddleware())
}
