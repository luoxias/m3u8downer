package interfaceList

import "io"

type Reader interface {
	GetReader() io.Reader
	Close() error
}
