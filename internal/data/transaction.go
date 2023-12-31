package data

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type Transaction struct {
	Id              int64           `json:"id"`
	TransactionDate TransactionDate `json:"transaction_date"`
	Amount          float64         `json:"amount"`
	UserID          int64           `json:"user_id"`
	CreatedAt       time.Time       `json:"created_at"`
}

type TransactionDate time.Time

func (td *TransactionDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(`"2006/01/02"`, s)
	if err != nil {
		return err
	}

	*td = TransactionDate(t)

	return nil
}

func (td TransactionDate) MarshalJSON() ([]byte, error) {
	jsonValue := td.Date().Format("2006/01/02")

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}

func (td TransactionDate) Date() time.Time {
	return time.Time(td)
}

type TransactionModel struct {
	DB *sql.DB
}

func (t TransactionModel) Insert(transaction *Transaction) error {
	query := `
		INSERT INTO transactions (txn_date, amount, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, txn_date, amount, user_id, created_at`

	args := []any{time.Time(transaction.TransactionDate), transaction.Amount, transaction.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.Id, &transaction.TransactionDate, &transaction.Amount, &transaction.UserID, &transaction.CreatedAt)
}

func (t TransactionModel) GetAll(id int64) ([]*Transaction, error) {
	query := `
		SELECT id, txn_date, amount, user_id, created_at FROM transactions t
		WHERE t.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.Id,
			&transaction.TransactionDate,
			&transaction.Amount,
			&transaction.UserID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t TransactionModel) GetAllBetweenDates(id int64, dateFrom, dateTo time.Time) ([]*Transaction, error) {
	query := `
		SELECT id, txn_date, amount, user_id, created_at FROM transactions t
		WHERE t.user_id = $1 AND t.txn_date >= $2 AND t.txn_date <= $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{id, dateFrom, dateTo}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.Id,
			&transaction.TransactionDate,
			&transaction.Amount,
			&transaction.UserID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t TransactionModel) InsertMultiple(transactions []*Transaction) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, transaction := range transactions {
		query := `
			INSERT INTO transactions (txn_date, amount, user_id)
			VALUES ($1, $2, $3)
			RETURNING id, txn_date, amount, user_id, created_at`

		args := []any{time.Time(transaction.TransactionDate), transaction.Amount, transaction.UserID}

		err = tx.QueryRowContext(ctx, query, args...).Scan(&transaction.Id, &transaction.TransactionDate, &transaction.Amount, &transaction.UserID, &transaction.CreatedAt)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
