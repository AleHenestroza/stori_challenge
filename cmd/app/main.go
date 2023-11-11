package main

import (
	"context"
	"database/sql"
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

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	smtp smtp
	dsn  string
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn, db, err := connectDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	cfg.dsn = dsn

	smtp, err := buildSmtpStruct()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	cfg.smtp = smtp
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
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func buildSmtpStruct() (smtp, error) {
	smtpHost, err := getEnv("SMTP_HOST")
	if err != nil {
		return smtp{}, err
	}
	smtpPortStr, err := getEnv("SMTP_PORT")
	if err != nil {
		return smtp{}, err
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return smtp{}, err
	}
	smtpUsername, err := getEnv("SMTP_USERNAME")
	if err != nil {
		return smtp{}, err
	}
	smtpPassword, err := getEnv("SMTP_PASSWORD")
	if err != nil {
		return smtp{}, err
	}
	smtpSender, err := getEnv("SMTP_SENDER")
	if err != nil {
		return smtp{}, err
	}

	smtp := smtp{
		host:     smtpHost,
		port:     smtpPort,
		username: smtpUsername,
		password: smtpPassword,
		sender:   smtpSender,
	}

	return smtp, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("could not read key %s", key)
}

func connectDB() (string, *sql.DB, error) {
	dsn, err := getEnv("STORI_DB_DSN")
	if err != nil {
		return "", nil, err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "", nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return "", nil, err
	}

	return dsn, db, nil
}
