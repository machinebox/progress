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
	N         int64
	Length    int64
	Remaining time.Duration
	Estimated time.Time
}

// Complete gets whether the operation is complete or not.
func (p Progress) Complete() bool {
	return p.N >= p.Length
}

// Percent calculates the percentage complete.
func (p Progress) Percent() float64 {
	n, length := float64(p.N), float64(p.Length)
	if n == 0 {
		return 0
	}
	if n == length {
		return 100
	}
	return 100 / (length / n)
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
				N:      n,
				Length: length,
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
				p.Estimated = started.Add(future)
				p.Remaining = p.Estimated.Sub(now)
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
