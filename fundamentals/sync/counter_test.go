package counter_test

import (
	counter "github.com/charlie-chiu/go_with_test/fundamentals/sync"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		c := counter.Counter{}
		c.Inc()
		c.Inc()
		c.Inc()

		assertCounter(t, &c, 3)
	})

	t.Run("it runs safely concurrently", func(t *testing.T) {
		wantedCount := 1000
		c := counter.Counter{}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func(w *sync.WaitGroup) {
				c.Inc()
				w.Done()
			}(&waitGroup)
		}
		waitGroup.Wait()

		assertCounter(t, &c, wantedCount)
	})
}

func assertCounter(t *testing.T, counter *counter.Counter, want int) {
	if counter.Value() != want {
		t.Errorf("got %d want %d", counter.Value(), want)
	}
}
