package transaction

import (
	"time"
)

type MonthlySummary struct {
	Month        string
	Balance      float32
	Transactions int
}

func NewMonthlySummary(txns []Transaction, month string, transactions int) (MonthlySummary, error) {
	balance, err := calculateMonthlyBalance(txns, month)
	if err != nil {
		return MonthlySummary{}, err
	}

	summary := MonthlySummary{
		Month:   month,
		Balance: balance,
		Transactions: transactions,
	}

	return summary, nil
}

func calculateMonthlyBalance(txns []Transaction, month string) (float32, error) {
	var totalBalance float32

	for _, t := range txns {
		date, err := t.GetDate()
		if err != nil {
			return 0, err
		}

		m, err := time.Parse("January", month)
		if err != nil {
			return 0, err
		}

		if date.Month() == m.Month() {
			totalBalance += t.Amount
		}
	}

	return totalBalance, nil
}
