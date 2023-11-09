package main

import (
	"fmt"
	"net/http"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
}

func (app *application) transactionsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	app.csvLoader.Read("./txns.csv")

	records, err := app.csvLoader.GetRecords()
	if err != nil {
		app.logger.Error("could not read records", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	transactions, err := app.parser.Parse(records)
	if err != nil {
		app.logger.Error("could not parse transactions", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	summary, err := transaction.NewMonthlySummary(transactions, "July")
	if err != nil {
		app.logger.Error("could not build monthly summary", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("Summary: %v", summary)
}
