package main

import (
	"fmt"
	"net/http"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

func (app *application) transactionsSummaryHandler(w http.ResponseWriter, r *http.Request) {
	app.csvLoader.Read("./txns.csv")

	records, err := app.csvLoader.GetRecords()
	if err != nil {
		fmt.Print(err)
	}

	transactions, err := app.parser.Parse(records)
	if err != nil {
		fmt.Print(err)
	}

	summary, err := transaction.NewMonthlySummary(transactions, "July")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Printf("Summary: %v", summary)
}
