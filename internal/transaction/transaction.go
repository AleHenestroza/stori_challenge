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
