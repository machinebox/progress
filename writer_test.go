package progress

import (
	"bytes"
	"io"
	"testing"

	"github.com/matryer/is"
)

func TestNewWriter(t *testing.T) {
	is := is.New(t)

	// check Writer interfaces
	var (
		_ io.Writer = (*Writer)(nil)
		_ Counter   = (*Writer)(nil)
	)

	var buf bytes.Buffer
	w := NewWriter(&buf, 5)

	n, err := w.Write([]byte("1"))
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(w.N(), int64(1)) // r.N()

	n, err = w.Write([]byte("1"))
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(w.N(), int64(2)) // r.N()

	n, err = w.Write([]byte("123"))
	is.NoErr(err)
	is.Equal(n, 3)            // n
	is.Equal(w.N(), int64(5)) // r.N()

}
