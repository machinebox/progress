# github.com/machinebox/progress

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

1. Wrap an `io.Reader` or `io.Writer` with `NewReader` and `NewWriter` respectively
1. You should specify the total number of bytes if known - otherwise, `Percent` and `Remaining` helpers will not work
1. Start a Goroutine to periodically check the progress, and do something with it - like log it
1. Use the readers and writers as normal
