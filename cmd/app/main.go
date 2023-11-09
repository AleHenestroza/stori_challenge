package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/reader"
	"github.com/alehenestroza/stori-backend-challenge/internal/parser"
)

type config struct {
	port int
	env  string
}

type application struct {
	config    config
	logger    *slog.Logger
	csvLoader reader.CsvDataReader
	parser    parser.TransactionParser
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config:    cfg,
		logger:    logger,
		csvLoader: *reader.NewCsvDataReader(),
		parser:    parser.NewTransactionParser(),
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
