package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/dal/mysql"
	"github.com/ahaostudy/calendar_reminder/job"
	"github.com/ahaostudy/calendar_reminder/log"
	"github.com/ahaostudy/calendar_reminder/middleware/ginmw"

	"github.com/ahaostudy/calendar_reminder/router"

	"github.com/joho/godotenv"

	"github.com/ahaostudy/calendar_reminder/conf"
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

	server := &http.Server{Addr: conf.GetConf().Server.Address, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal(err)
		}
	}()

	// async jobs
	job.InitAsyncJobs()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// waiting for remaining requests to be processed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}

	shutdown()
}

func initMiddleware(r *gin.Engine) {
	r.Use(ginmw.LoggerMiddleware())
	r.Use(ginmw.CorsMiddleware())
	r.Use(ginmw.GlobalAuthMiddleware())
}

// close all resources
func shutdown() {
	job.DestroyAsyncJobs()
}
