package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
)

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

func (app *application) listTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	txns, err := app.models.Transactions.GetAll(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": txns}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) saveTransactionHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	var input struct {
		TransactionDate string  `json:"transaction_date"`
		Amount          float64 `json:"amount"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	date, err := time.Parse("2006/01/02", input.TransactionDate)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid date format"))
		return
	}

	transaction := data.Transaction{
		TransactionDate: data.TransactionDate(date),
		Amount:          input.Amount,
		UserID:          user.ID,
	}

	err = app.models.Transactions.Insert(&transaction)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
