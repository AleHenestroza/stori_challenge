package main

import (
	"net/http"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     "1.0.0",
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) transactionsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	type emailFields struct {
		AccountName         string
		TotalBalance        string
		AccountSummary      transaction.AccountSummary
		AverageDebitAmount  string
		AverageCreditAmount string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	records, err := app.csvLoader.Read("./txns.csv")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	transactions, err := app.parser.Parse(records)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	summary := transaction.NewAccountSummary(transactions)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.mailer.Send(input.Email, "account_summary.tmpl", emailFields{
		AccountName:         input.Name,
		TotalBalance:        summary.Balance,
		AccountSummary:      summary,
		AverageDebitAmount:  summary.DebitAverage,
		AverageCreditAmount: summary.CreditAverage,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	w.WriteHeader(http.StatusAccepted)
}
