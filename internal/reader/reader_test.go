package reader_test

import (
	"os"
	"testing"

	"github.com/alehenestroza/stori-backend-challenge/internal/reader"
)

func TestCsvDataReader_ReadFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	data := "header1,header2\nvalue1,value2\nvalue3,value4"
	if _, err := tmpFile.Write([]byte(data)); err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("error closing temp file: %v", err)
	}

	csvReader := reader.NewCsvDataReader()

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("error opening temp file: %v", err)
	}
	defer file.Close()

	records, err := csvReader.ReadFile(file)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	expected := []string{"header1,header2", "value1,value2", "value3,value4"}
	if !compareSlices(records, expected) {
		t.Errorf("got %v but expected %v", records, expected)
	}
}

func compareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
