package data

import (
	"context"
	"database/sql"
	"time"
)

type Transaction struct {
	Id        int64           `json:"id"`
	Date      TransactionDate `json:"date"`
	Amount    float64         `json:"amount"`
	UserID    int64           `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
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

	args := []any{transaction.Date, transaction.Amount, transaction.UserID}

	return t.DB.QueryRow(query, args...).Scan(&transaction.Id, &transaction.Date, &transaction.Amount, &transaction.UserID, &transaction.CreatedAt)
}

func (t TransactionModel) GetAll(id int64) ([]*Transaction, error) {
	query := `
		SELECT id, txn_date, amount, user_id FROM transactions t
		WHERE t.user_id = $1`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*Transaction{}

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.Id,
			&transaction.Date,
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
