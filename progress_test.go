package progress

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestTicker(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	c := &counter{}
	var size int64 = 200
	ticker := NewTicker(ctx, c, size, 10*time.Millisecond)
	var events []Progress
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case tick, ok := <-ticker:
				if !ok {
					return
				}
				events = append(events, tick)
			}
		}
	}()
	go func() {
		for {
			n := c.N() + 50
			c.SetN(n)
			if n >= size {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	// wait for it to finish
	<-done
	is.True(len(events) > 5) // should be >5 events
	is.Equal(events[len(events)-1].Complete(), true)
}

func TestProgress(t *testing.T) {
	is := is.New(t)

	now := time.Now()

	is.Equal((Progress{n: 1}).N(), int64(1))
	is.Equal((Progress{estimated: now}).Estimated(), now)
	is.Equal((Progress{remaining: 1 * time.Minute}).Remaining(), 1*time.Minute)
	is.Equal((Progress{size: 10}).Size(), int64(10))
	is.Equal((Progress{n: 1, size: 2}).Complete(), false)
	is.Equal((Progress{n: 2, size: 2}).Complete(), true)
	is.Equal((Progress{n: 0, size: 2}).Percent(), 0.0)
	is.Equal((Progress{n: 1, size: 2}).Percent(), 50.0)
	is.Equal((Progress{n: 2, size: 2}).Percent(), 100.0)

}

type counter struct {
	n int64
}

func (c *counter) N() int64 {
	return atomic.LoadInt64(&c.n)
}

func (c *counter) SetN(n int64) {
	atomic.StoreInt64(&c.n, n)
}
