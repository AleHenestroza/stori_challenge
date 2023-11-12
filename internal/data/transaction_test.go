package data

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestUnmarshalTransaction(t *testing.T) {
	jsonData := `{"id": 1, "date": "2023/11/12", "amount": 100.0}`

	var transaction Transaction
	if err := json.Unmarshal([]byte(jsonData), &transaction); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	expectedDate := time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC)
	if time.Time(transaction.Date) != expectedDate {
		t.Errorf("expected %v but got %v", expectedDate, transaction.Date)
	}
}
