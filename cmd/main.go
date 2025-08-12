package main

import (
	"io"
	"os"

	"github.com/financial_tracer/internal/config"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	log := NewLogger(cfg.App.Env)

	_, err := postgresql.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: router(gin)
	//TODO: handler(gin)

}

func NewLogger(cfg string) *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("logu.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mx := io.MultiWriter(os.Stdout, file)
	log.Out = mx

	switch cfg {
	case "Local":
		log.SetLevel(logrus.DebugLevel)

	case "Debug":
		log.SetLevel(logrus.DebugLevel)

	case "prod":
		log.SetLevel(logrus.InfoLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
