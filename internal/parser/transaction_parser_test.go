package parser_test

import (
	"testing"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
	"github.com/alehenestroza/stori-backend-challenge/internal/parser"
)

func TestParse(t *testing.T) {
	p := parser.NewTransactionParser()
	rows := []string{
		"1,01/09,100.50",
		"2,11/1,50.25",
	}

	transactions, err := p.Parse(rows)
	if err != nil {
		t.Errorf("Parse failed with error: %v", err)
	}

	expectedTransactions := []data.Transaction{
		{Id: 1, Date: data.TransactionDate(time.Date(0, 1, 9, 0, 0, 0, 0, time.UTC)), Amount: 100.50},
		{Id: 2, Date: data.TransactionDate(time.Date(0, 11, 1, 0, 0, 0, 0, time.UTC)), Amount: 50.25},
	}
	if !compareTransactions(transactions, expectedTransactions) {
		t.Errorf("Transactions mismatch. Got: %v, Expected: %v", transactions, expectedTransactions)
	}
}

func compareTransactions(a, b []data.Transaction) bool {
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
