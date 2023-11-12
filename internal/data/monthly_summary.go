package data

import (
	"time"
)

type MonthlySummary struct {
	Month        time.Month
	Balance      float64
	Transactions int
}

func NewMonthlySummary(txns []Transaction, month time.Month, transactions int) (MonthlySummary, error) {
	balance, err := calculateMonthlyBalance(txns, month)
	if err != nil {
		return MonthlySummary{}, err
	}

	summary := MonthlySummary{
		Month:        month,
		Balance:      balance,
		Transactions: transactions,
	}

	return summary, nil
}

func calculateMonthlyBalance(txns []Transaction, month time.Month) (float64, error) {
	var totalBalance float64

	for _, t := range txns {
		if t.TransactionDate.Date().Month() == month {
			totalBalance += t.Amount
		}
	}

	return totalBalance, nil
}
