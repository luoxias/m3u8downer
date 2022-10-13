package implements

import (
	"io"
	"os"
)

type ReadFromFile struct {
	file   *os.File
	reader io.Reader
}

func (r *ReadFromFile) SetFileName(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	r.file = file
	r.reader = file
	return nil
}

func (r *ReadFromFile) GetReader() io.Reader {
	return r.reader
}

func (r *ReadFromFile) Close() error {
	return r.file.Close()
}
