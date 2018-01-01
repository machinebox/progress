# progress

`io.Reader` and `io.Writer` with progress and remaining time estimation.

## Usage

```go
import (
	"github.com/machinebox/progress"
)

// get a reader and the total expected number of bytes
s := `Now that's what I call progress`
size := len(s)
r := progress.NewReader(strings.NewReader(s), int64(size))

// Start a goroutine printing progress
go func(){
	defer log.Printf("done")
	interval := 1 * time.Second
	ticker := progress.NewTicker(r, interval)
	for {
		select {
		case tick, ok := <-ticker:
			if !ok {
				// done
				return
			}
			log.Printf("%f%% completed, about %v remaining", tick.Percent(). tick.Remaining())
		}
	}
}()

// use the Reader as normal
if _, err := io.Copy(dest, r); err != nil {
	log.Fatalln(err)
}
```

