package progress

import (
	"io"
	"sync/atomic"
)

// Reader counts the bytes read through it.
type Reader struct {
	r io.Reader
	n int64
}

// NewReader makes a new Reader that counts the bytes
// read through it.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	atomic.AddInt64(&r.n, int64(n))
	return
}

// N gets the number of bytes that have been read
// so far.
func (r *Reader) N() int64 {
	return atomic.LoadInt64(&r.n)
}
