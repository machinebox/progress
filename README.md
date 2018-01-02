# github.com/machinebox/progress

`io.Reader` and `io.Writer` with progress and remaining time estimation.

## Usage

```go
ctx := context.Background()

// get a reader and the total expected number of bytes
s := `Now that's what I call progress`
size := len(s)
r := progress.NewReader(strings.NewReader(s))

// Start a goroutine printing progress
go func(){
	defer log.Printf("done")
	interval := 1 * time.Second
	progressChan := progress.NewTicker(ctx, r, size, interval)
	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				// if ok is false, the process is finished
				return
			}
			log.Printf("about %v remaining...", progress.Remaining())
		}
	}
}()

// use the Reader as normal
if _, err := io.Copy(dest, r); err != nil {
	log.Fatalln(err)
}
```

1. Wrap an `io.Reader` or `io.Writer` with `NewReader` and `NewWriter` respectively
1. Capture the total number of expeted bytes
1. Use `progress.NewTicker` to get a channel on which progress updates will be sent
1. Start a Goroutine to periodically check the progress, and do something with it - like log it
1. Use the readers and writers as normal
