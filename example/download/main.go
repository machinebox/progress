package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/machinebox/progress"
	"github.com/pkg/errors"
)

func main() {
	if err := run(os.Args[1:]...); err != nil {
		log.Fatalln(err)
	}
}

func run(args ...string) error {
	if len(args) < 1 {
		return errors.New("bad number of arguments")
	}
	url := args[0]
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()
	contentLengthHeader := resp.Header.Get("Content-Length")
	if contentLengthHeader == "" {
		return errors.New("cannot determine progress without Content-Length")
	}
	size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "bad Content-Length %q", contentLengthHeader)
	}
	ctx := context.Background()
	r := progress.NewReader(resp.Body)
	go func() {
		progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)
		for p := range progressChan {
			fmt.Printf("\r%v remaining...", p.Remaining().Round(time.Second))
		}
		fmt.Println("\rdownload is completed")
	}()
	if _, err := io.Copy(ioutil.Discard, r); err != nil {
		return errors.Wrap(err, "failed to read body")
	}
	return nil
}
