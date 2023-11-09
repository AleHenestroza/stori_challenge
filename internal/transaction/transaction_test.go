package transaction_test

import (
	"testing"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

func TestTransactionDateParsing(t *testing.T) {
	testCases := []struct {
		inputDate   string
		expected    time.Time
		expectError bool
	}{
		{
			inputDate:   "7/15",
			expected:    time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			inputDate:   "invalid",
			expected:    time.Time{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.inputDate, func(t *testing.T) {
			transaction := transaction.Transaction{Date: tc.inputDate}
			date, err := transaction.GetDate()

			if tc.expectError {
				if err == nil {
					t.Errorf("an error was expected, but none was produced")
				}
			} else {
				if err != nil {
					t.Errorf("an error was produced, but none was expected")
				}

				if date != tc.expected {
					t.Errorf("incorrect date. expected: %v but got: %v", tc.expected, date)
				}
			}
		})
	}
}
