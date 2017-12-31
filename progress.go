package progress

import (
	"time"
)

// Counter counts bytes.
// Both Reader and Writer are Counter types.
type Counter interface {
	// N is the number of bytes that have been read
	// so far.
	N() int64
	// Len is the total number of bytes expected.
	Len() int64
}

// Percent calculates the percentage complete for the
// specified Counter.
func Percent(c Counter) float64 {
	n, length := float64(c.N()), float64(c.Len())
	if n == 0 {
		return 0
	}
	if n == length {
		return 100
	}
	return 100 / (length / n)
}

// Complete gets whether the Counter is complete or not.
func Complete(c Counter) bool {
	return c.N() >= c.Len()
}

// Progress represents a moment of progress.
type Progress struct {
	n         int64
	length    int64
	remaining time.Duration
	estimated time.Time
}

// Complete gets whether the operation is complete or not.
func (p Progress) Complete() bool {
	return p.n >= p.length
}

// Percent calculates the percentage complete.
func (p Progress) Percent() float64 {
	n, length := float64(p.n), float64(p.length)
	if n == 0 {
		return 0
	}
	if n == length {
		return 100
	}
	return 100 / (length / n)
}

// Remaining gets the amount of time until the operation is
// expected to be finished. Use Estimated to get a fixed completion time.
func (p Progress) Remaining() time.Duration {
	return p.remaining
}

// Estimated gets the time at which the operation is expected
// to finish. Use Reamining to get a Duration.
func (p Progress) Estimated() time.Time {
	return p.estimated
}

// NewTicker gets a channel on which ticks of Progress are sent
// at duration d intervals until the operation is complete.
func NewTicker(c Counter, d time.Duration) <-chan Progress {
	var started time.Time
	ch := make(chan Progress)
	go func() {
		defer close(ch)
		for {
			n, length := c.N(), c.Len()
			p := Progress{
				n:      n,
				length: length,
			}
			nF, lengthF := float64(n), float64(length)
			if started.IsZero() {
				if n > 0 {
					started = time.Now()
				}
			} else {
				now := time.Now()
				ratio := nF / lengthF
				past := now.Sub(started)
				future := time.Duration(float64(past) / ratio)
				p.estimated = started.Add(future)
				p.remaining = p.estimated.Sub(now)
			}
			ch <- p
			if n >= length {
				return
			}
			time.Sleep(d)
		}
	}()
	return ch
}
