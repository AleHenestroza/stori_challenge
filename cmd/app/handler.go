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
	// TODO: Refactor this handler
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

	summary := transaction.NewAccountSummary(transactions)
	if err != nil {
		app.logger.Error("could not build monthly summary", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	accountSummary := app.formater.FormatTransactions(summary.MonthlySummary)

	err = app.mailer.Send("alehenestroza@gmail.com", "account_summary.tmpl", struct {
		TotalBalance        string
		AccountSummary      []string
		AverageDebitAmount  string
		AverageCreditAmount string
	}{
		TotalBalance:        summary.Balance,
		AccountSummary:      accountSummary,
		AverageDebitAmount:  summary.DebitAverage,
		AverageCreditAmount: summary.CreditAverage,
	})
	if err != nil {
		app.logger.Error("could not send email", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("Summary: %v", summary)
}
