package log

import (
	"io"
	"os"

	"github.com/ahaostudy/calendar_reminder/conf"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	logFile, err := os.OpenFile(conf.GetConf().Server.LogPath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		logrus.Fatalf("open log file err: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
