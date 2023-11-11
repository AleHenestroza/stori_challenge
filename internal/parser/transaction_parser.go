package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alehenestroza/stori-backend-challenge/internal/transaction"
)

type TransactionParser struct{}

func NewTransactionParser() TransactionParser {
	return TransactionParser{}
}

func (tp TransactionParser) Parse(rows []string) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, len(rows))

	for i, row := range rows {
		id, date, amount, err := tp.parseRow(row)
		if err != nil {
			fmt.Printf("could not parse line %d: %s. Error: %v", i, row, err)
		}

		transactions[i] = transaction.Transaction{
			Id:     int32(id),
			Date:   tp.formatDate(date),
			Amount: float32(amount),
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

func (tp TransactionParser) formatDate(date string) string {
	parts := strings.Split(date, "/")
	parts[0] = fmt.Sprintf("%02s", parts[0])
	parts[1] = fmt.Sprintf("%02s", parts[1])
	return strings.Join(parts, "/")
}
