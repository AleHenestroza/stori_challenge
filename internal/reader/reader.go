package reader

import "os"

type DataReader interface {
	ReadFile(file *os.File) ([]string, error)
}
