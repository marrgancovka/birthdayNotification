package main

import (
	_ "birthdayNotification/docs"
	"birthdayNotification/internal/app"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

// @title Birthday Notification API
// @version 1.0
// @description This is a server for birthday notifications.

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name jwt-birthday-service

// @host localhost:8080
// @BasePath /api
func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	application := app.NewApp(logger)

	if err := application.Start(); err != nil {
		logger.Fatal(err)
	}

}
