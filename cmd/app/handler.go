package main

import (
	"net/http"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
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
		AccountSummary      data.AccountSummary
		AverageDebitAmount  string
		AverageCreditAmount string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	records, err := app.csvLoader.Read("./txns.csv")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	transactions, err := app.parser.Parse(records)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	summary := data.NewAccountSummary(transactions)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *application) sendTransactionSummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if user.ID != id {
		app.notPermittedResponse(w, r)
		return
	}

	txns, err := app.models.Transactions.GetAll(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": txns}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
