package mailer

import (
	"fmt"
	"strconv"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

type Formater struct{}

func NewFormater() Formater {
	return Formater{}
}

func (f Formater) FormatTransactions(txns []transaction.MonthlySummary) []string {
	var listItems []string

	for _, summary := range txns {
		transactionsStr := strconv.Itoa(summary.Transactions)

		li := fmt.Sprintf("Number of transactions in %s: %s", summary.Month, transactionsStr)

		listItems = append(listItems, li)
	}

	return listItems
}
