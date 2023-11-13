package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
)

const csvFile = "./txns.csv"

type emailFields struct {
	AccountName         string
	TotalBalance        string
	AccountSummary      data.AccountSummary
	AverageDebitAmount  string
	AverageCreditAmount string
}

func (app *application) sendLocalTransactionsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	file, err := os.Open(csvFile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer file.Close()

	transactions, err := app.parser.Parse(file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	summary, err := data.NewAccountSummary(transactions)
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

func (app *application) saveTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	file, err := app.getRequestFile(w, r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	transactions, err := app.parser.Parse(file)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	for _, t := range transactions {
		t.UserID = user.ID
		err = app.models.Transactions.Insert(t)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) sendTransactionsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	transactions, err := app.models.Transactions.GetAll(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	accountSummary, err := data.NewAccountSummary(transactions)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.mailer.Send(user.Email, "account_summary.tmpl", emailFields{
		AccountName:         user.Name,
		TotalBalance:        accountSummary.Balance,
		AccountSummary:      accountSummary,
		AverageDebitAmount:  accountSummary.DebitAverage,
		AverageCreditAmount: accountSummary.CreditAverage,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
