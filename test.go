package test

import "errors"

// ErrPlain is the default error that is returned for functions in this package.
var ErrPlain = errors.New("error")

////////////////

// ErrorReader implements an io.Reader that will do N successive reads before it returns ErrPlain.
type ErrorReader struct {
	n int
}

// NewErrorReader returns a new ErrorReader.
func NewErrorReader(n int) *ErrorReader {
	return &ErrorReader{n}
}

// Read implements the io.Reader interface.
func (r *ErrorReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.n == 0 {
		return 0, ErrPlain
	}
	r.n--
	b[0] = '.'
	return 1, nil
}

////////////////

// ErrorWriter implements an io.Writer that will do N successive writes before it returns ErrPlain.
type ErrorWriter struct {
	n int
}

// NewErrorWriter returns a new ErrorWriter.
func NewErrorWriter(n int) *ErrorWriter {
	return &ErrorWriter{n}
}

// Write implements the io.Writer interface.
func (w *ErrorWriter) Write(b []byte) (n int, err error) {
	if w.n == 0 {
		return 0, ErrPlain
	}
	w.n--
	return len(b), nil
}
