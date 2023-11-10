package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/mailer"
	"github.com/alehenestroza/stori-backend-challenge/internal/parser"
	"github.com/alehenestroza/stori-backend-challenge/internal/reader"
)

type config struct {
	port int
	env  string
	smtp smtp
}

type smtp struct {
	host     string
	port     int
	username string
	password string
	sender   string
}

type application struct {
	config    config
	logger    *slog.Logger
	csvLoader reader.CsvDataReader
	parser    parser.TransactionParser
	mailer    mailer.Mailer
	formater  mailer.Formater
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "<username>", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "<password>", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Stori Test <no-reply@storitest.com>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	formater := mailer.NewFormater()
	mailer := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)

	app := &application{
		config:    cfg,
		logger:    logger,
		csvLoader: *reader.NewCsvDataReader(),
		parser:    parser.NewTransactionParser(),
		mailer:    mailer,
		formater:  formater,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
