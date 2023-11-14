package reader

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"strings"
)

type CsvDataReader struct{}

func NewCsvDataReader() CsvDataReader {
	return CsvDataReader{}
}

func (r CsvDataReader) ReadFile(file multipart.File) ([]string, error) {
	reader := csv.NewReader(file)
	var rows []string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		rows = append(rows, strings.Join(row, ","))
	}

	return rows, nil
}
