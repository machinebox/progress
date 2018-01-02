package progress

import (
	"io"
	"sync/atomic"
)

// Writer counts the bytes written.
type Writer struct {
	w io.Writer
	n int64
}

// NewWriter gets a Writer that counts the number
// of bytes written.
// The length should be set to the expected number of bytes to be written.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	atomic.AddInt64(&w.n, int64(n))
	return
}

// N gets the number of bytes that have been written
// so far.
func (w *Writer) N() int64 {
	return atomic.LoadInt64(&w.n)
}
