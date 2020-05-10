package packet

import "io"

type ReadWriter interface {
	Reader
	Writer
}

type readWriter struct {
	Reader
	Writer
}

func NewReadWriter(rw io.ReadWriter) ReadWriter {
	return &readWriter{
		Reader: NewReader(rw),
		Writer: NewWriter(rw),
	}
}
