package reader

import (
	"os"
	"strings"
)

type CsvDataReader struct{}

func NewCsvDataReader() CsvDataReader {
	return CsvDataReader{}
}

func (r CsvDataReader) ReadFile(file *os.File) ([]string, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	text := string(buf)
	rows := strings.Split(text, "\n")
	return rows, nil
}
