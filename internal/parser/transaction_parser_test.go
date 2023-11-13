package parser_test

import (
	"os"
	"testing"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
	"github.com/alehenestroza/stori-backend-challenge/internal/parser"
	"github.com/alehenestroza/stori-backend-challenge/internal/reader"
)

func TestParse(t *testing.T) {
	p := parser.NewTransactionParser(reader.NewCsvDataReader())

	tmpFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	fileData := "Id,Date,Transaction\n1,01/09,100.50\n2,11/1,50.25"
	if _, err := tmpFile.Write([]byte(fileData)); err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("error closing temp file: %v", err)
	}

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("error opening temp file: %v", err)
	}
	defer file.Close()

	transactions, err := p.Parse(file)
	if err != nil {
		t.Errorf("Parse failed with error: %v", err)
	}

	expectedTransactions := []*data.Transaction{
		{
			Id:              1,
			TransactionDate: data.TransactionDate(time.Date(0, 1, 9, 0, 0, 0, 0, time.UTC)),
			Amount:          100.50,
		},
		{
			Id:              2,
			TransactionDate: data.TransactionDate(time.Date(0, 11, 1, 0, 0, 0, 0, time.UTC)),
			Amount:          50.25,
		},
	}
	if !compareTransactions(transactions, expectedTransactions) {
		t.Errorf("Transactions mismatch. Got: %v, Expected: %v", transactions, expectedTransactions)
	}
}

func compareTransactions(a, b []*data.Transaction) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Id != b[i].Id {
			return false
		}

		if a[i].TransactionDate != b[i].TransactionDate {
			return false
		}

		if a[i].Amount != b[i].Amount {
			return false
		}
	}
	return true
}
