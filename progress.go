package progress

import (
	"context"
	"time"
)

// Counter counts bytes.
// Both Reader and Writer are Counter types.
type Counter interface {
	// N is the number of bytes that have been read
	// or written so far.
	N() int64
}

// Progress represents a moment of progress.
type Progress struct {
	n         int64
	size      int64
	remaining time.Duration
	estimated time.Time
}

// N gets the total number of bytes read or written
// so far.
func (p Progress) N() int64 {
	return p.n
}

// Size gets the total number of bytes that are expected to
// be read or written.
func (p Progress) Size() int64 {
	return p.size
}

// Complete gets whether the operation is complete or not.
func (p Progress) Complete() bool {
	return p.n >= p.size
}

// Percent calculates the percentage complete.
func (p Progress) Percent() float64 {
	n, length := float64(p.n), float64(p.size)
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
// The size is the total number of expected bytes being read or written.
// Cancellable via context.
func NewTicker(ctx context.Context, counter Counter, size int64, d time.Duration) <-chan Progress {
	var (
		started time.Time
		ch      = make(chan Progress)
	)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
				progress := Progress{
					n:    counter.N(),
					size: size,
				}
				if started.IsZero() {
					if progress.n > 0 {
						started = time.Now()
					}
				} else {
					now := time.Now()
					nF, lengthF := float64(progress.n), float64(size)
					ratio := nF / lengthF
					past := now.Sub(started)
					future := time.Duration(float64(past) / ratio)
					progress.estimated = started.Add(future)
					progress.remaining = progress.estimated.Sub(now)
				}
				ch <- progress
				if progress.Complete() {
					return
				}
			}
		}
	}()
	return ch
}
