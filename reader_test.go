package progress

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestNewReader(t *testing.T) {
	is := is.New(t)

	// check Reader interfaces
	var (
		_ io.Reader = (*Reader)(nil)
		_ Counter   = (*Reader)(nil)
	)

	s := `Now that's what I call progress`
	r := NewReader(strings.NewReader(s))

	buf := make([]byte, 1)
	n, err := r.Read(buf)
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(r.N(), int64(1)) // r.N()

	n, err = r.Read(buf)
	is.NoErr(err)
	is.Equal(n, 1)            // n
	is.Equal(r.N(), int64(2)) // r.N()

	// read to the end
	b, err := ioutil.ReadAll(r)
	is.NoErr(err)
	is.Equal(len(b), 29)       // len(b)
	is.Equal(r.N(), int64(31)) // r.N()

}
