package data

import (
	"fmt"
	"slices"
	"time"
)

type AccountSummary struct {
	MonthlySummary []MonthlySummary
	Balance        string
	DebitAverage   string
	CreditAverage  string
}

func NewAccountSummary(transactions []*Transaction) (AccountSummary, error) {
	var balance, debitAmount, creditAmount float64
	var debits, credits int
	summary := AccountSummary{}
	monthMap := make(map[time.Month][]*Transaction)

	for _, t := range transactions {
		month := t.TransactionDate.Date().Month()

		if _, found := monthMap[month]; !found {
			monthMap[month] = []*Transaction{}
		}

		balance += t.Amount

		if t.Amount < 0 {
			debits += 1
			debitAmount += t.Amount
		} else {
			credits += 1
			creditAmount += t.Amount
		}

		monthMap[month] = append(monthMap[month], t)
	}

	for month := range monthMap {
		monthlySummary, err := NewMonthlySummary(monthMap[month], month, len(monthMap[month]))
		if err != nil {
			return AccountSummary{}, err
		}

		summary.MonthlySummary = append(summary.MonthlySummary, monthlySummary)
	}

	slices.SortFunc(summary.MonthlySummary, func(a, b MonthlySummary) int {
		if a.Month > b.Month {
			return 1
		}

		if a.Month < b.Month {
			return -1
		}

		return 0
	})

	summary.Balance = fmt.Sprintf("%.2f", balance)
	if credits > 0 {
		summary.CreditAverage = fmt.Sprintf("%.2f", creditAmount/float64(credits))
	} else {
		summary.CreditAverage = "0.00"
	}

	if debits > 0 {
		summary.DebitAverage = fmt.Sprintf("%.2f", debitAmount/float64(debits))
	} else {
		summary.DebitAverage = "0.00"
	}

	return summary, nil
}
