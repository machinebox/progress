# `progress` [![GoDoc](https://godoc.org/github.com/machinebox/progress?status.png)](http://godoc.org/github.com/machinebox/progress) [![Build Status](https://travis-ci.org/machinebox/progress.svg?branch=master)](https://travis-ci.org/machinebox/progress) [![Go Report Card](https://goreportcard.com/badge/github.com/machinebox/progress)](https://goreportcard.com/report/github.com/machinebox/progress)

`io.Reader` and `io.Writer` with progress and remaining time estimation.

## Usage

```go
ctx := context.Background()

// get a reader and the total expected number of bytes
s := `Now that's what I call progress`
size := len(s)
r := progress.NewReader(strings.NewReader(s))

// Start a goroutine printing progress
go func() {
	ctx := context.Background()
	progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)
	for p := range progressChan {
		fmt.Printf("\r%v remaining...", p.Remaining().Round(time.Second))
	}
	fmt.Println("\rdownload is completed")
}()

// use the Reader as normal
if _, err := io.Copy(dest, r); err != nil {
	log.Fatalln(err)
}
```

1. Wrap an `io.Reader` or `io.Writer` with `NewReader` and `NewWriter` respectively
1. Capture the total number of expected bytes
1. Use `progress.NewTicker` to get a channel on which progress updates will be sent
1. Start a Goroutine to periodically check the progress, and do something with it - like log it
1. Use the readers and writers as normal
