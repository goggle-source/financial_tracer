package main

import (
	"io"
	"net/http"
	"os"

	"github.com/financial_tracer/internal/config"
	"github.com/financial_tracer/internal/handlers"
	categoryHandlers "github.com/financial_tracer/internal/handlers/categories"
	userHandlers "github.com/financial_tracer/internal/handlers/user"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/financial_tracer/internal/servic/category"
	"github.com/financial_tracer/internal/servic/user"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	log := NewLogger(cfg.App.Env)

	db, err := postgresql.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	users := user.CreateUserServer(db)
	handlersUser := userHandlers.CreateHandlersUser(cfg.App.SercretKey, users, log)
	categories := category.CreateCategoryServer(db)
	handlersCategory := categoryHandlers.CreateHandlersCategory(categories, log)
	r := handlers.Router(handlersUser, handlersCategory, log, cfg.App.SercretKey)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  cfg.Server.IdleTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
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
