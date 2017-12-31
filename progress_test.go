package progress

import (
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestTicker(t *testing.T) {
	is := is.New(t)
	c := &counter{0, 200}
	ticker := NewTicker(c, 10*time.Millisecond)
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
			if n >= c.Len() {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	// wait for it to finish
	<-done
	is.True(len(events) > 5) // should be >5 events
}

func TestProgress(t *testing.T) {
	is := is.New(t)

	is.Equal((Progress{n: 1, length: 2}).Complete(), false)
	is.Equal((Progress{n: 2, length: 2}).Complete(), true)

	is.Equal((Progress{n: 0, length: 2}).Percent(), 0.0)
	is.Equal((Progress{n: 1, length: 2}).Percent(), 50.0)
	is.Equal((Progress{n: 2, length: 2}).Percent(), 100.0)

}

func XTestTickerTimes(t *testing.T) {
	is := is.New(t)
	c := &counter{0, 200}
	ticker := NewTicker(c, 100*time.Millisecond)
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
				log.Printf("%+v", tick)
			}
		}
	}()
	go func() {
		for {
			n := c.N() + 1
			c.SetN(n)
			if n >= c.Len() {
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	// wait for it to finish
	<-done
	is.True(len(events) > 5) // should be >5 events
}

func TestPercent(t *testing.T) {
	is := is.New(t)

	c := &counter{0, 200}
	is.Equal(int(Percent(c)), 0)
	c.n = 50
	is.Equal(int(Percent(c)), 25)
	c.n = 100
	is.Equal(int(Percent(c)), 50)
	c.n = 150
	is.Equal(int(Percent(c)), 75)
	c.n = 200
	is.Equal(int(Percent(c)), 100)

}

func TestComplete(t *testing.T) {
	is := is.New(t)

	c := &counter{0, 200}
	is.Equal(Complete(c), false)
	c.n = 50
	is.Equal(Complete(c), false)
	c.n = 100
	is.Equal(Complete(c), false)
	c.n = 150
	is.Equal(Complete(c), false)
	c.n = 200
	is.Equal(Complete(c), true)
}

type counter struct {
	n, l int64
}

func (c *counter) N() int64 {
	return atomic.LoadInt64(&c.n)
}

func (c *counter) SetN(n int64) {
	atomic.StoreInt64(&c.n, n)
}

func (c *counter) Len() int64 {
	return atomic.LoadInt64(&c.l)
}
