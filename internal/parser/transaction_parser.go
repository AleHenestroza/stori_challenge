package parser

import (
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alehenestroza/stori-backend-challenge/internal/reader"

	"github.com/alehenestroza/stori-backend-challenge/internal/data"
)

type TransactionParser struct {
	Reader reader.DataReader
}

func NewTransactionParser(reader reader.DataReader) TransactionParser {
	return TransactionParser{Reader: reader}
}

func (tp TransactionParser) Parse(file multipart.File) ([]*data.Transaction, error) {
	rows, err := tp.Reader.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var transactions []*data.Transaction

	for i, row := range rows {
		if isHeaderRow(row) {
			continue
		}

		if !isValidDataRow(row) {
			return nil, fmt.Errorf("invalid character at row %d: %s", i+1, row)
		}

		id, dateStr, amount, err := tp.parseRow(row)
		if err != nil {
			return nil, err
		}

		date, err := tp.parseDate(dateStr)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &data.Transaction{
			Id:              id,
			TransactionDate: data.TransactionDate(date),
			Amount:          amount,
		})
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

func isValidDataRow(row string) bool {
	parts := strings.Split(row, ",")
	regexPattern := "^[0-9/.+\\-]+$"
	re := regexp.MustCompile(regexPattern)

	for _, part := range parts {
		if !re.MatchString(part) {
			return false
		}
	}

	return true
}
