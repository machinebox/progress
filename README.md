# progress

`io.Reader` and `io.Writer` with progress and remaining time estimation.

## Usage

```go
// get a reader and the total expected number of bytes
s := `Now that's what I call progress`
size := len(s)
r := NewReader(strings.NewReader(s), int64(size))

// Start a goroutine printing progress
go func(){
	ticker := progress.NewTicker(r, 1 * time.Second)
	for {
		select {
		case tick, ok := <-ticker:
			if !ok {
				// done
				return
			}
			log.Printf("%f%% remaining, about %v", tick.Percent(). tick.Remaining)
		}
	}
}()

// use the Reader as normal
if _, err := io.Copy(dest, r); err != nil {
	log.Fatalln(err)
}
```

