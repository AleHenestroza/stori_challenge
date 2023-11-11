package mailer

import (
	"fmt"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

type Formater struct{}

func NewFormater() Formater {
	return Formater{}
}

func (f Formater) FormatTransactions(txns []transaction.MonthlySummary) map[string]int {
	items := make(map[string]int, 12)

	for _, summary := range txns {
		items[fmt.Sprintf("Number of transactions in %s", summary.Month)] = summary.Transactions
	}

	return items
}
