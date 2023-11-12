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
		id, dateStr, amount, err := tp.parseRow(row)
		if err != nil {
			return nil, err
		}

		date, err := tp.parseDate(dateStr)
		if err != nil {
			return nil, err
		}

		transactions[i] = data.Transaction{
			Id:     int64(id),
			Date:   data.TransactionDate(date),
			Amount: float64(amount),
		}
	}

	return transactions, nil
}

func (tp TransactionParser) parseRow(row string) (int32, string, float32, error) {
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

	return int32(id), date, float32(amount), nil
}

func (tp TransactionParser) parseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("1/2", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
