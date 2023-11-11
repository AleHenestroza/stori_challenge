package transaction

import (
	"time"
)

const dateFormat = "1/2"

type Transaction struct {
	Id     int32
	Date   string
	Amount float32
}

func (t Transaction) GetDate() (time.Time, error) {
	date, err := time.Parse(dateFormat, t.Date)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func CompareTransactionsByDate(t1, t2 Transaction) (int, error) {
	date1, err := t1.GetDate()
	if err != nil {
		return 0, err
	}
	date2, err := t2.GetDate()
	if err != nil {
		return 0, err
	}

	if date1.Before(date2) {
		return -1, nil
	} else if date1.After(date2) {
		return 1, nil
	} else {
		return 0, nil
	}
}
