package transaction

import (
	"time"
)

type MonthlySummary struct {
	Month        time.Month
	Balance      float32
	Transactions int
}

func NewMonthlySummary(txns []Transaction, month time.Month, transactions int) (MonthlySummary, error) {
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

func calculateMonthlyBalance(txns []Transaction, month time.Month) (float32, error) {
	var totalBalance float32

	for _, t := range txns {
		date, err := t.GetDate()
		if err != nil {
			return 0, err
		}

		if date.Month() == month {
			totalBalance += t.Amount
		}
	}

	return totalBalance, nil
}
