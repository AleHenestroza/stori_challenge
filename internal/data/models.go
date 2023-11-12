package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Transactions TransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Transactions: TransactionModel{db},
	}
}
