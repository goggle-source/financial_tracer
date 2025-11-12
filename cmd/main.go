package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/financial_tracer/internal/config"
	"github.com/financial_tracer/internal/handlers"
	categoryHandlers "github.com/financial_tracer/internal/handlers/categories"
	transactionHandlers "github.com/financial_tracer/internal/handlers/transaction"
	userHandlers "github.com/financial_tracer/internal/handlers/user"
	"github.com/financial_tracer/internal/infastructure/cash"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/financial_tracer/internal/servic/category"
	"github.com/financial_tracer/internal/servic/transaction"
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

	red := cash.CreateRealRedis(*cfg)
	ctx := context.Background()

	users := user.CreateUserServer(db, db, db, cfg.App.SercretKey, log)
	handlersUser := userHandlers.CreateHandlersUser(cfg.App.SercretKey, users, users, users, log, ctx)
	categories := category.CreateCategoryServer(db, db, db, db, db, log, &red)
	handlersCategory := categoryHandlers.CreateHandlersCategory(categories, categories, categories, categories, categories, log, ctx)
	transactions := transaction.CreateTransactionServer(db, db, db, db, log, &red)
	handlersTransaction := transactionHandlers.CreateTransactionHandlers(transactions, transactions, transactions, transactions, log, ctx)
	r := handlers.Router(handlersUser, handlersCategory, log, handlersTransaction, cfg.App.SercretKey)

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
