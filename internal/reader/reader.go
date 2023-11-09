package reader

import (
	"errors"
	"os"
	"strings"
)

type CsvDataReader struct {
	Records []string
}

func NewCsvDataReader() *CsvDataReader {
	return &CsvDataReader{}
}

func (csvl *CsvDataReader) Read(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	rawString := string(file)
	csvl.Records = strings.Split(rawString, "\n")[1:] // Skip first line, which contains the headers of the CSV file

	return nil
}

func (csvl *CsvDataReader) GetRecords() ([]string, error) {
	if len(csvl.Records) == 0 {
		return nil, errors.New("no records present")
	}
	return csvl.Records, nil
}
