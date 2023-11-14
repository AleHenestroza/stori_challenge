package reader

import "mime/multipart"

type DataReader interface {
	ReadFile(file multipart.File) ([]string, error)
}
