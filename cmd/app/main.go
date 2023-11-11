package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.Parse()

	cfg.smtp = buildSmtpStruct()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mailer := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)

	app := &application{
		config:    cfg,
		logger:    logger,
		csvLoader: *reader.NewCsvDataReader(),
		parser:    parser.NewTransactionParser(),
		mailer:    mailer,
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

func buildSmtpStruct() smtp {
	smtpHost, err := getEnv("SMTP_HOST")
	if err != nil {
		panic(err)
	}
	smtpPortStr, err := getEnv("SMTP_PORT")
	if err != nil {
		panic(err)
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		panic(err)
	}
	smtpUsername, err := getEnv("SMTP_USERNAME")
	if err != nil {
		panic(err)
	}
	smtpPassword, err := getEnv("SMTP_PASSWORD")
	if err != nil {
		panic(err)
	}
	smtpSender, err := getEnv("SMTP_SENDER")
	if err != nil {
		panic(err)
	}

	smtp := smtp{
		host:     smtpHost,
		port:     smtpPort,
		username: smtpUsername,
		password: smtpPassword,
		sender:   smtpSender,
	}

	return smtp
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("could not read key %s", key)
}
