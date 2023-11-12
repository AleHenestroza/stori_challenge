package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
)

type TransactionParser struct{}

func NewTransactionParser() TransactionParser {
	return TransactionParser{}
}

func (tp TransactionParser) Parse(rows []string) ([]data.Transaction, error) {
	transactions := make([]data.Transaction, len(rows))

	for i, row := range rows {
		if isHeaderRow(row) {
			continue
		}

		id, dateStr, amount, err := tp.parseRow(row)
		if err != nil {
			return nil, err
		}

		date, err := tp.parseDate(dateStr)
		if err != nil {
			return nil, err
		}

		transactions[i] = data.Transaction{
			Id:              id,
			TransactionDate: data.TransactionDate(date),
			Amount:          amount,
		}
	}

	return transactions, nil
}

func (tp TransactionParser) parseRow(row string) (int64, string, float64, error) {
	parts := strings.Split(row, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("wrong input format")
	}

	id, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("could not parse transaction ID")
	}

	date := parts[1]

	amount, err := strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("could not parse transaction Amount")
	}
	
	return id, date, amount, nil
}

func (tp TransactionParser) parseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("1/2", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func isHeaderRow(row string) bool {
	parts := strings.Split(row, ",")

	return strings.EqualFold(parts[0], "Id") && strings.EqualFold(parts[1], "Date") && strings.EqualFold(parts[2], "Transaction")
}
