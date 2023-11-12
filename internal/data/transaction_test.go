package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalTransaction(t *testing.T) {
	jsonData := `{"id": 1, "transaction_date": "2023/11/12", "amount": 100.0}`

	var transaction Transaction
	if err := json.Unmarshal([]byte(jsonData), &transaction); err != nil {
		t.Error("Error decoding JSON:", err)
	}

	expectedDate := time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC)
	if time.Time(transaction.TransactionDate) != expectedDate {
		t.Errorf("expected %v but got %v", expectedDate, transaction.TransactionDate)
	}
}
