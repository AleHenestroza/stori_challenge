package parser_test

import (
	"testing"

	"github.com/alehenestroza/stori-backend-challenge/internal/parser"
	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

func TestParse(t *testing.T) {
	p := parser.NewTransactionParser()
	rows := []string{
		"1,2023-11-09,100.50",
		"2,2023-11-10,50.25",
	}

	transactions, err := p.Parse(rows)
	if err != nil {
		t.Errorf("Parse failed with error: %v", err)
	}

	expectedTransactions := []transaction.Transaction{
		{Id: 1, Date: "2023-11-09", Amount: 100.50},
		{Id: 2, Date: "2023-11-10", Amount: 50.25},
	}
	if !compareTransactions(transactions, expectedTransactions) {
		t.Errorf("Transactions mismatch. Got: %v, Expected: %v", transactions, expectedTransactions)
	}
}

func compareTransactions(a, b []transaction.Transaction) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
